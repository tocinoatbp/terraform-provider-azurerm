package redisenterprise_test

import (
	"context"
	"fmt"
	"testing"

	"github.com/hashicorp/go-azure-helpers/lang/response"
	"github.com/hashicorp/terraform-provider-azurerm/internal/services/redisenterprise/sdk/2021-08-01/databases"

	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance"
	"github.com/hashicorp/terraform-provider-azurerm/internal/acceptance/check"
	"github.com/hashicorp/terraform-provider-azurerm/internal/clients"
	"github.com/hashicorp/terraform-provider-azurerm/internal/tf/pluginsdk"
	"github.com/hashicorp/terraform-provider-azurerm/utils"
)

type RedisenterpriseDatabaseResource struct{}

func TestRedisEnterpriseDatabase_basic(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func TestRedisEnterpriseDatabase_requiresImport(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.basic(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.RequiresImportErrorStep(r.requiresImport),
	})
}

func TestRedisEnterpriseDatabase_complete(t *testing.T) {
	data := acceptance.BuildTestData(t, "azurerm_redis_enterprise_database", "test")
	r := RedisenterpriseDatabaseResource{}
	data.ResourceTest(t, r, []acceptance.TestStep{
		{
			Config: r.complete(data),
			Check: acceptance.ComposeTestCheckFunc(
				check.That(data.ResourceName).ExistsInAzure(r),
			),
		},
		data.ImportStep(),
	})
}

func (r RedisenterpriseDatabaseResource) Exists(ctx context.Context, client *clients.Client, state *pluginsdk.InstanceState) (*bool, error) {
	id, err := databases.ParseDatabaseID(state.ID)
	if err != nil {
		return nil, err
	}

	resp, err := client.RedisEnterprise.DatabaseClient.Get(ctx, *id)
	if err != nil {
		if response.WasNotFound(resp.HttpResponse) {
			return utils.Bool(false), nil
		}

		return nil, fmt.Errorf("retrieving %s: %+v", *id, err)
	}

	return utils.Bool(true), nil
}

func (r RedisenterpriseDatabaseResource) template(data acceptance.TestData) string {
	// I have to hardcode the location because some features are not currently available in all regions
	return fmt.Sprintf(`
provider "azurerm" {
  features {}
}

resource "azurerm_resource_group" "test" {
  name     = "acctestRG-redisEnterprise-%d"
  location = "%s"
}

resource "azurerm_redis_enterprise_cluster" "test" {
  name                = "acctest-rec-%d"
  resource_group_name = azurerm_resource_group.test.name
  location            = azurerm_resource_group.test.location

  sku_name = "Enterprise_E20-4"
}
`, data.RandomInteger, "eastus", data.RandomInteger)
}

func (r RedisenterpriseDatabaseResource) basic(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "test" {
  name                = "default"
  resource_group_name = azurerm_resource_group.test.name
  cluster_id          = azurerm_redis_enterprise_cluster.test.id
}
`, template)
}

func (r RedisenterpriseDatabaseResource) requiresImport(data acceptance.TestData) string {
	config := r.basic(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "import" {
  name                = azurerm_redis_enterprise_database.test.name
  resource_group_name = azurerm_redis_enterprise_database.test.resource_group_name
  cluster_id          = azurerm_redis_enterprise_database.test.cluster_id
}
`, config)
}

func (r RedisenterpriseDatabaseResource) complete(data acceptance.TestData) string {
	template := r.template(data)
	return fmt.Sprintf(`
%s

resource "azurerm_redis_enterprise_database" "test" {
  resource_group_name = azurerm_resource_group.test.name
  cluster_id          = azurerm_redis_enterprise_cluster.test.id

  client_protocol   = "Encrypted"
  clustering_policy = "EnterpriseCluster"
  eviction_policy   = "NoEviction"

  module {
    name = "RediSearch"
    args = ""
  }

  module {
    name = "RedisBloom"
    args = "ERROR_RATE 0.00 INITIAL_SIZE 400"
  }

  module {
    name = "RedisTimeSeries"
    args = "RETENTION_POLICY 20"
  }

  port = 10000
}
`, template)
}
