package es

import (
	"context"
	"errors"

	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/schema"

	elastic7 "github.com/olivere/elastic/v7"
	elastic6 "gopkg.in/olivere/elastic.v6"
)

func dataSourceElasticsearchAlias() *schema.Resource {
	return &schema.Resource{
		Description: "`elasticsearch_alias` can be used to retrieve alias for the provider's current elasticsearch cluster.",
		Read:        dataSourceElasticsearchAliasRead,

		Schema: map[string]*schema.Schema{
			"name": {
				Type:        schema.TypeString,
				Required:    true,
				Description: "Name of the alias to retrieve",
			},
			"exists": {
				Type:        schema.TypeBool,
				Computed:    true,
				Description: "should be set to `true` if alias exists",
			},
			"alias": {
				Type:        schema.TypeMap,
				Computed:    true,
				Description: "alias metadata in active elasticsearch cluster",
			},
		},
	}
}

func dataSourceElasticsearchAliasRead(d *schema.ResourceData, m interface{}) error {
	// The upstream elastic client does not export the property for the urls
	// it's using. Presumably the URLS would be available where the client is
	// intantiated, but in terraform, that's not always practicable.

	/*
		HEAD _alias/logs_write-1
		{"statusCode":404,"error":"Not Found","message":"404 - Not Found"}
		HEAD _alias/logs_write
		200 - OK

		GET _alias/logs_write
		{
			"logs-000001" : {
				"aliases" : {
				"logs_write" : { }
				}
			}
		}
	*/
	aliasName := d.Get("name").(string)

	var err error
	esClient, err := getClient(m.(*ProviderConf))
	if err != nil {
		return err
	}

	var url string
	switch client := esClient.(type) {
	case *elastic7.Client:
		exists, alias, err = aliasElasticsearchGet(client, aliasName)
		// urls := reflect.ValueOf(client).Elem().FieldByName("urls")
		// if urls.Len() > 0 {
		// 	url = urls.Index(0).String()
		// }
	case *elastic6.Client:
		// urls := reflect.ValueOf(client).Elem().FieldByName("urls")
		// if urls.Len() > 0 {
		// 	url = urls.Index(0).String()
		// }
	default:
		return errors.New("this version of Elasticsearch is not supported")
	}
	d.SetId(url)
	err = d.Set("url", url)

	return err
}

func aliasElasticsearchGet(client *elastic7.Client, name string) (bool, map[string]interface{}, error) {
	result, err := client.Aliases().Alias(name).Do(context.TODO())
	if err != nil {
		return false, nil, err
	}
	for indiceName, indiceResult := range result.Indices {
		return indiceResult.HasAlias(name), indiceResult, nil
	}

}
