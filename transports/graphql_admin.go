package transports

import (
	"context"
	"log"

	"github.com/BuxOrg/bux"
	"github.com/machinebox/graphql"
	"github.com/mrz1836/go-datastore"
)

// NewXpub will register an xPub
func (g *TransportGraphQL) NewXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {

	// adding a xpub needs to be signed by an admin key
	if g.adminXPriv == nil {
		return ErrAdminKey
	}

	reqBody := `
   	mutation ($metadata: Metadata) {
	  xpub(
		xpub: "` + rawXPub + `"
		metadata: $metadata
	  ) {
	    id
	  }
	}`
	req := graphql.NewRequest(reqBody)
	req.Var(FieldMetadata, processMetadata(metadata))
	variables := map[string]interface{}{
		FieldMetadata: processMetadata(metadata),
	}

	bodyString, err := getBodyString(reqBody, variables)
	if err != nil {
		return err
	}
	if err = addSignature(&req.Header, g.adminXPriv, bodyString); err != nil {
		return err
	}

	// run it and capture the response
	var xPubData interface{}

	return g.client.Run(ctx, req, &xPubData)
}

// RegisterXpub alias for NewXpub
func (g *TransportGraphQL) RegisterXpub(ctx context.Context, rawXPub string, metadata *bux.Metadata) error {
	return g.NewXpub(ctx, rawXPub, metadata)
}

// AdminGetStatus get whether admin key is valid
func (g *TransportGraphQL) AdminGetStatus(ctx context.Context) (bool, error) {

	reqBody := `
	query {
	  admin_get_status
	}`

	var status bool
	if err := g.doGraphQLAdminQuery(ctx, reqBody, nil, &status); err != nil {
		return false, err
	}

	return status, nil
}

// AdminGetStats get admin stats
func (g *TransportGraphQL) AdminGetStats(ctx context.Context) (*bux.AdminStats, error) {

	reqBody := `
	  query {
        admin_get_stats {
          balance
          destinations
          transactions
          paymails
          utxos
          xpubs
          transactions_per_day
          utxos_per_type
        }
      }`

	var stats *bux.AdminStats
	if err := g.doGraphQLAdminQuery(ctx, reqBody, nil, &stats); err != nil {
		return nil, err
	}

	return stats, nil
}

// AdminGetAccessKeys get all access keys filtered by conditions
func (g *TransportGraphQL) AdminGetAccessKeys(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.AccessKey, error) {

	var models []*bux.AccessKey
	method := `admin_access_keys_list`
	fields := `
      id
      xpub_id
      key
      metadata
      created_at
      updated_at
      deleted_at
      revoked_at
    `

	if err := g.adminGetModels(ctx, conditions, metadata, queryParams, method, fields, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetAccessKeysCount get a count of all the access keys filtered by conditions
func (g *TransportGraphQL) AdminGetAccessKeysCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return g.adminCount(ctx, conditions, metadata, "admin_access_keys_count")
}

// AdminGetBlockHeaders get all block headers filtered by conditions
func (g *TransportGraphQL) AdminGetBlockHeaders(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.BlockHeader, error) {

	var models []*bux.BlockHeader
	method := `admin_block_headers_list`
	fields := `
	  id
	  height
	  time
	  nonce
	  version
	  hash_previous_block
	  hash_merkle_root
	  bits
	  synced
	  created_at
	  updated_at
	  deleted_at
    `

	if err := g.adminGetModels(ctx, conditions, metadata, queryParams, method, fields, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetBlockHeadersCount get a count of all the block headers filtered by conditions
func (g *TransportGraphQL) AdminGetBlockHeadersCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return g.adminCount(ctx, conditions, metadata, "admin_block_headers_count")
}

// AdminGetDestinations get all block destinations filtered by conditions
func (g *TransportGraphQL) AdminGetDestinations(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Destination, error) {

	var models []*bux.Destination
	method := `admin_destinations_list`
	fields := `
	  id
	  xpub_id
	  locking_script
	  type
	  chain
	  num
	  address
	  draft_id
	  metadata
	  created_at
	  updated_at
	  deleted_at
    `

	if err := g.adminGetModels(ctx, conditions, metadata, queryParams, method, fields, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetDestinationsCount get a count of all the destinations filtered by conditions
func (g *TransportGraphQL) AdminGetDestinationsCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return g.adminCount(ctx, conditions, metadata, "admin_destinations_count")
}

// AdminGetPaymail get a paymail by address
func (g *TransportGraphQL) AdminGetPaymail(ctx context.Context, address string) (*bux.PaymailAddress, error) {

	reqBody := `
	  query ($address: String!) {
        admin_paymail_get (
          address: $address
        ) {
          id
          xpub_id
          alias
          domain
          public_name
          avatar
          created_at
          updated_at
          deleted_at
        }
      }`

	variables := map[string]interface{}{
		FieldAddress: address,
	}

	var paymail *bux.PaymailAddress
	if err := g.doGraphQLAdminQuery(ctx, reqBody, variables, &paymail); err != nil {
		return nil, err
	}

	return paymail, nil
}

// AdminGetPaymails get all block paymails filtered by conditions
func (g *TransportGraphQL) AdminGetPaymails(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.PaymailAddress, error) {

	var models []*bux.PaymailAddress
	method := `admin_paymails_list`
	fields := `
	  id
	  xpub_id
	  alias
	  domain
	  public_name
	  avatar
	  created_at
	  updated_at
	  deleted_at
    `

	if err := g.adminGetModels(ctx, conditions, metadata, queryParams, method, fields, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetPaymailsCount get a count of all the paymails filtered by conditions
func (g *TransportGraphQL) AdminGetPaymailsCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return g.adminCount(ctx, conditions, metadata, "admin_paymails_count")
}

// AdminCreatePaymail create a new paymail for a xpub
func (g *TransportGraphQL) AdminCreatePaymail(ctx context.Context, xPubID string, address string, publicName string,
	avatar string) (*bux.PaymailAddress, error) {

	reqBody := `
      mutation (
        $xpub_id: String!
        $address: String!
        $public_name: String!
        $avatar: String!
      ) {
        admin_paymail_create (
          xpub_id: $xpub_id
          address: $address
          public_name: $public_name
          avatar: $avatar
        ) {
          id
          xpub_id
          alias
          domain
          public_name
          avatar
          created_at
          updated_at
          deleted_at
        }
      }`

	variables := map[string]interface{}{
		FieldXpubID:     xPubID,
		FieldAddress:    address,
		FieldPublicName: publicName,
		FieldAvatar:     avatar,
	}

	var paymail *bux.PaymailAddress
	if err := g.doGraphQLAdminQuery(ctx, reqBody, variables, &paymail); err != nil {
		return nil, err
	}

	return paymail, nil
}

// AdminDeletePaymail delete a paymail address from the database
func (g *TransportGraphQL) AdminDeletePaymail(ctx context.Context, address string) (*bux.PaymailAddress, error) {

	reqBody := `
      mutation (
        $address: String!
      ) {
        admin_paymail_delete (
          address: $address
        ) {
          id
          xpub_id
          alias
          domain
          public_name
          avatar
          created_at
          updated_at
          deleted_at
        }
      }`

	variables := map[string]interface{}{
		FieldAddress: address,
	}

	var paymail *bux.PaymailAddress
	if err := g.doGraphQLAdminQuery(ctx, reqBody, variables, &paymail); err != nil {
		return nil, err
	}

	return paymail, nil
}

// AdminGetTransactions get all block transactions filtered by conditions
func (g *TransportGraphQL) AdminGetTransactions(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Transaction, error) {

	var models []*bux.Transaction
	method := `admin_transactions_list`
	fields := `
	  id
	  hex
	  block_hash
	  block_height
	  fee
	  number_of_inputs
	  number_of_outputs
	  output_value
	  total_value
	  metadata
	  created_at
	  updated_at
	  deleted_at
    `

	if err := g.adminGetModels(ctx, conditions, metadata, queryParams, method, fields, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetTransactionsCount get a count of all the transactions filtered by conditions
func (g *TransportGraphQL) AdminGetTransactionsCount(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata) (int64, error) {

	return g.adminCount(ctx, conditions, metadata, "admin_transactions_count")
}

// AdminGetUtxos get all block utxos filtered by conditions
func (g *TransportGraphQL) AdminGetUtxos(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Utxo, error) {

	var models []*bux.Utxo
	method := `admin_utxos_list`
	fields := `
	  id
	  xpub_id
	  satoshis
	  script_pub_key
	  type
	  draft_id
	  reserved_at
	  spending_tx_id
	  created_at
	  updated_at
	  deleted_at
    `

	if err := g.adminGetModels(ctx, conditions, metadata, queryParams, method, fields, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetUtxosCount get a count of all the utxos filtered by conditions
func (g *TransportGraphQL) AdminGetUtxosCount(ctx context.Context, conditions map[string]interface{}, metadata *bux.Metadata) (int64, error) {
	return g.adminCount(ctx, conditions, metadata, "admin_utxos_count")
}

// AdminGetXPubs get all block xpubs filtered by conditions
func (g *TransportGraphQL) AdminGetXPubs(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams) ([]*bux.Xpub, error) {

	var models []*bux.Xpub
	method := `admin_xpubs_list`
	fields := `
	  id
	  current_balance
	  next_internal_num
	  next_external_num
	  metadata
	  created_at
	  updated_at
	  deleted_at
    `

	if err := g.adminGetModels(ctx, conditions, metadata, queryParams, method, fields, &models); err != nil {
		return nil, err
	}

	return models, nil
}

// AdminGetXPubsCount get a count of all the xpubs filtered by conditions
func (g *TransportGraphQL) AdminGetXPubsCount(ctx context.Context, conditions map[string]interface{}, metadata *bux.Metadata) (int64, error) {
	return g.adminCount(ctx, conditions, metadata, "admin_xpubs_count")
}

func (g *TransportGraphQL) adminGetModels(ctx context.Context, conditions map[string]interface{},
	metadata *bux.Metadata, queryParams *datastore.QueryParams, method string, fields string,
	models interface{}) error {

	reqBody := `
	  query ($conditions: Map, $metadata: Metadata, $params: QueryParams) {
        ` + method + ` (
          conditions: $conditions
          metadata: $metadata
          params: $params
        ) {
          ` + fields + `
        }
      }`

	variables := map[string]interface{}{
		FieldConditions:  conditions,
		FieldMetadata:    processMetadata(metadata),
		FieldQueryParams: queryParams,
	}

	return g.doGraphQLAdminQuery(ctx, reqBody, variables, &models)
}

func (g *TransportGraphQL) adminCount(ctx context.Context, conditions map[string]interface{}, metadata *bux.Metadata,
	method string) (int64, error) {

	// adding a xpub needs to be signed by an admin key
	if g.adminXPriv == nil {
		return 0, ErrAdminKey
	}

	reqBody := `
   	   query ($conditions: Map, $metadata: Metadata) {
        ` + method + ` (
          conditions: $conditions
          metadata: $metadata
        )
      }`

	req := graphql.NewRequest(reqBody)
	req.Var(FieldConditions, conditions)
	req.Var(FieldMetadata, processMetadata(metadata))
	variables := map[string]interface{}{
		FieldConditions: conditions,
		FieldMetadata:   processMetadata(metadata),
	}

	bodyString, err := getBodyString(reqBody, variables)
	if err != nil {
		return 0, err
	}
	if err = addSignature(&req.Header, g.adminXPriv, bodyString); err != nil {
		return 0, err
	}

	// run it and capture the response
	var count int64
	if err = g.client.Run(ctx, req, &count); err != nil {
		return 0, err
	}

	return count, nil
}

func (g *TransportGraphQL) doGraphQLAdminQuery(ctx context.Context, reqBody string, variables map[string]interface{},
	respData interface{}) error {

	req := graphql.NewRequest(reqBody)
	for key, value := range variables {
		req.Var(key, value)
	}

	err := g.signGraphQLRequest(req, reqBody, variables, g.adminXPriv, nil)
	if err != nil {
		return err
	}

	// run it and capture the response
	if err = g.client.Run(ctx, req, &respData); err != nil {
		return err
	}
	if g.debug {
		log.Printf("model: %v\n", respData)
	}

	return nil
}

// AdminRecordTransaction will record a transaction as an admin
func (g *TransportGraphQL) AdminRecordTransaction(ctx context.Context, hex string) (*bux.Transaction, error) {

	reqBody := `
   	mutation() {
	  admin_transaction (
		hex:"` + hex + `",
	  ) {
		id
		hex
		block_hash
		block_height
		fee
		number_of_inputs
		number_of_outputs
		output_value
		total_value
		direction
		metadata
		created_at
		updated_at
		deleted_at
	  }
	}`
	req := graphql.NewRequest(reqBody)

	err := g.signGraphQLRequest(req, reqBody, nil, g.xPriv, g.xPub)
	if err != nil {
		return nil, err
	}

	// run it and capture the response
	var respData NewTransactionData
	if err = g.client.Run(ctx, req, &respData); err != nil {
		return nil, err
	}
	transaction := respData.Transaction
	if g.debug {
		log.Printf("transaction: %s\n", transaction.ID)
	}

	return transaction, nil
}
