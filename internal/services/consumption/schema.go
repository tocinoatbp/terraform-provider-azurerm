package consumption

import (
	"github.com/Azure/azure-sdk-for-go/services/consumption/mgmt/2019-10-01/consumption"
	"github.com/hashicorp/terraform-provider-azurerm/helpers/azure"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/consumption/validate"
	validateManagementGroup "github.com/hashicorp/terraform-provider-azurerm/internal/services/managementgroup/validate"
	resourceValidate "github.com/hashicorp/terraform-provider-azurerm/internal/services/resource/validate"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/validation"
)

func SchemaConsumptionBudgetResourceGroupResource() map[string]*pluginsdk.Schema {
	resourceGroupNameSchema := map[string]*pluginsdk.Schema{
		"resource_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: resourceValidate.ResourceGroupID,
		},
	}

	return azure.MergeSchema(SchemaConsumptionBudgetCommonResource(), resourceGroupNameSchema)
}

func SchemaConsumptionBudgetSubscriptionResource() map[string]*pluginsdk.Schema {
	subscriptionIDSchema := map[string]*pluginsdk.Schema{
		"subscription_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validation.IsUUID,
		},
	}

	return azure.MergeSchema(SchemaConsumptionBudgetCommonResource(), subscriptionIDSchema)
}

func SchemaConsumptionBudgetFilterDimensionElement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					"ChargeType",
					"Frequency",
					"InvoiceId",
					"Meter",
					"MeterCategory",
					"MeterSubCategory",
					"PartNumber",
					"PricingModel",
					"Product",
					"ProductOrderId",
					"ProductOrderName",
					"PublisherType",
					"ReservationId",
					"ReservationName",
					"ResourceGroupName",
					"ResourceGuid",
					"ResourceId",
					"ResourceLocation",
					"ResourceType",
					"ServiceFamily",
					"ServiceName",
					"UnitOfMeasure",
				}, false),
			},
			"operator": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "In",
				ValidateFunc: validation.StringInSlice([]string{
					"In",
				}, false),
			},
			"values": {
				Type:     pluginsdk.TypeList,
				MinItems: 1,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaConsumptionBudgetFilterTagElement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"name": {
				Type:     pluginsdk.TypeString,
				Required: true,
			},
			"operator": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  "In",
				ValidateFunc: validation.StringInSlice([]string{
					"In",
				}, false),
			},
			"values": {
				Type:     pluginsdk.TypeList,
				Required: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaConsumptionBudgetNotificationElement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
			"threshold": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			// Issue: https://github.com/Azure/azure-rest-api-specs/issues/16240
			// Toggling between these two values doesn't work at the moment and also doesn't throw an error
			// but it seems unlikely that a user would switch the threshold_type of their budgets frequently
			"threshold_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(consumption.ThresholdTypeActual),
				ForceNew: true, // todo: remove this when the above issue is fixed
				ValidateFunc: validation.StringInSlice([]string{
					string(consumption.ThresholdTypeActual),
					"Forecasted",
				}, false),
			},
			"operator": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(consumption.OperatorTypeEqualTo),
					string(consumption.OperatorTypeGreaterThan),
					string(consumption.OperatorTypeGreaterThanOrEqualTo),
				}, false),
			},

			"contact_emails": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"contact_groups": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},

			"contact_roles": {
				Type:     pluginsdk.TypeList,
				Optional: true,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}

func SchemaConsumptionBudgetCommonResource() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConsumptionBudgetName(),
		},

		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},

		"amount": {
			Type:         pluginsdk.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(1.0),
		},

		"filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dimension": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterDimensionElement()),
						Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
					},
					"tag": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterTagElement()),
						Elem:         SchemaConsumptionBudgetFilterTagElement(),
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
					},
					"not": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dimension": {
									Type:         pluginsdk.TypeList,
									MaxItems:     1,
									Optional:     true,
									ExactlyOneOf: []string{"filter.0.not.0.tag"},
									Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
								},
								"tag": {
									Type:         pluginsdk.TypeList,
									MaxItems:     1,
									Optional:     true,
									ExactlyOneOf: []string{"filter.0.not.0.dimension"},
									Elem:         SchemaConsumptionBudgetFilterTagElement(),
								},
							},
						},
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
					},
				},
			},
		},

		"notification": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 5,
			Set:      pluginsdk.HashResource(SchemaConsumptionBudgetNotificationElement()),
			Elem:     SchemaConsumptionBudgetNotificationElement(),
		},

		"time_grain": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(consumption.TimeGrainTypeMonthly),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(consumption.TimeGrainTypeBillingAnnual),
				string(consumption.TimeGrainTypeBillingMonth),
				string(consumption.TimeGrainTypeBillingQuarter),
				string(consumption.TimeGrainTypeAnnually),
				string(consumption.TimeGrainTypeMonthly),
				string(consumption.TimeGrainTypeQuarterly),
			}, false),
		},

		"time_period": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_date": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ConsumptionBudgetTimePeriodStartDate,
						ForceNew:     true,
					},
					"end_date": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
	}
}

func SchemaConsumptionBudgetManagementGroupResource() map[string]*pluginsdk.Schema {
	return map[string]*pluginsdk.Schema{
		"name": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validate.ConsumptionBudgetName(),
		},

		"management_group_id": {
			Type:         pluginsdk.TypeString,
			Required:     true,
			ForceNew:     true,
			ValidateFunc: validateManagementGroup.ManagementGroupID,
		},

		"etag": {
			Type:     pluginsdk.TypeString,
			Computed: true,
			Optional: true,
		},

		"amount": {
			Type:         pluginsdk.TypeFloat,
			Required:     true,
			ValidateFunc: validation.FloatAtLeast(1.0),
		},

		"filter": {
			Type:     pluginsdk.TypeList,
			Optional: true,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"dimension": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterDimensionElement()),
						Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
					},
					"tag": {
						Type:         pluginsdk.TypeSet,
						Optional:     true,
						Set:          pluginsdk.HashResource(SchemaConsumptionBudgetFilterTagElement()),
						Elem:         SchemaConsumptionBudgetFilterTagElement(),
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
					},
					"not": {
						Type:     pluginsdk.TypeList,
						Optional: true,
						MaxItems: 1,
						Elem: &pluginsdk.Resource{
							Schema: map[string]*pluginsdk.Schema{
								"dimension": {
									Type:         pluginsdk.TypeList,
									MaxItems:     1,
									Optional:     true,
									ExactlyOneOf: []string{"filter.0.not.0.tag"},
									Elem:         SchemaConsumptionBudgetFilterDimensionElement(),
								},
								"tag": {
									Type:         pluginsdk.TypeList,
									MaxItems:     1,
									Optional:     true,
									ExactlyOneOf: []string{"filter.0.not.0.dimension"},
									Elem:         SchemaConsumptionBudgetFilterTagElement(),
								},
							},
						},
						AtLeastOneOf: []string{"filter.0.dimension", "filter.0.tag", "filter.0.not"},
					},
				},
			},
		},

		"notification": {
			Type:     pluginsdk.TypeSet,
			Required: true,
			MinItems: 1,
			MaxItems: 5,
			Set:      pluginsdk.HashResource(SchemaConsumptionBudgetNotificationManagementGroupElement()),
			Elem:     SchemaConsumptionBudgetNotificationManagementGroupElement(),
		},

		"time_grain": {
			Type:     pluginsdk.TypeString,
			Optional: true,
			Default:  string(consumption.TimeGrainTypeMonthly),
			ForceNew: true,
			ValidateFunc: validation.StringInSlice([]string{
				string(consumption.TimeGrainTypeBillingAnnual),
				string(consumption.TimeGrainTypeBillingMonth),
				string(consumption.TimeGrainTypeBillingQuarter),
				string(consumption.TimeGrainTypeAnnually),
				string(consumption.TimeGrainTypeMonthly),
				string(consumption.TimeGrainTypeQuarterly),
			}, false),
		},

		"time_period": {
			Type:     pluginsdk.TypeList,
			Required: true,
			MinItems: 1,
			MaxItems: 1,
			Elem: &pluginsdk.Resource{
				Schema: map[string]*pluginsdk.Schema{
					"start_date": {
						Type:         pluginsdk.TypeString,
						Required:     true,
						ValidateFunc: validate.ConsumptionBudgetTimePeriodStartDate,
						ForceNew:     true,
					},
					"end_date": {
						Type:         pluginsdk.TypeString,
						Optional:     true,
						Computed:     true,
						ValidateFunc: validation.IsRFC3339Time,
					},
				},
			},
		},
	}
}

func SchemaConsumptionBudgetNotificationManagementGroupElement() *pluginsdk.Resource {
	return &pluginsdk.Resource{
		Schema: map[string]*pluginsdk.Schema{
			"enabled": {
				Type:     pluginsdk.TypeBool,
				Optional: true,
				Default:  true,
			},
			"threshold": {
				Type:         pluginsdk.TypeInt,
				Required:     true,
				ValidateFunc: validation.IntBetween(0, 1000),
			},
			// Issue: https://github.com/Azure/azure-rest-api-specs/issues/16240
			// Toggling between these two values doesn't work at the moment and also doesn't throw an error
			// but it seems unlikely that a user would switch the threshold_type of their budgets frequently
			"threshold_type": {
				Type:     pluginsdk.TypeString,
				Optional: true,
				Default:  string(consumption.ThresholdTypeActual),
				ForceNew: true, // todo: remove this when the above issue is fixed
				ValidateFunc: validation.StringInSlice([]string{
					string(consumption.ThresholdTypeActual),
					"Forecasted",
				}, false),
			},
			"operator": {
				Type:     pluginsdk.TypeString,
				Required: true,
				ValidateFunc: validation.StringInSlice([]string{
					string(consumption.OperatorTypeEqualTo),
					string(consumption.OperatorTypeGreaterThan),
					string(consumption.OperatorTypeGreaterThanOrEqualTo),
				}, false),
			},

			"contact_emails": {
				Type:     pluginsdk.TypeList,
				Required: true,
				MinItems: 1,
				Elem: &pluginsdk.Schema{
					Type:         pluginsdk.TypeString,
					ValidateFunc: validation.StringIsNotEmpty,
				},
			},
		},
	}
}
