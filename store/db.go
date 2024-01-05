package store

import (
	"authonomy/models"
	"authonomy/pkg/utils"
	"bytes"
	"encoding/json"

	"github.com/dgraph-io/badger/v3"
)

const (
	app_prefix             = "app-"
	policy_prefix          = "policy-"
	issued_policy_prefix   = "issued-"
	auth_prefix            = "auth-"
	conf_prefix            = "conf-"
	provider_schema_prefix = "prov-"
)

// Store encapsulates the BadgerDB operations
type Store struct {
	db     *badger.DB
	secret []byte
}

// NewStore initializes and returns a new Store instance
func NewStore(path string, s string) (*Store, error) {
	db, err := badger.Open(badger.DefaultOptions(path))
	if err != nil {
		return nil, err
	}
	secret := utils.GenerateEncryptionKey(s)
	return &Store{db: db, secret: secret}, nil
}

// ClearDB deletes all key-value pairs in the database
func (s *Store) ClearDB() error {
	return s.db.DropAll()
}

// Close safely closes the BadgerDB instance
func (s *Store) Close() {
	s.db.Close()
}

// SetApp stores an Application instance in the database
func (s *Store) SetApp(app models.ApplicationResponse) error {
	return s.db.Update(func(txn *badger.Txn) error {
		appJSON, err := json.Marshal(app)
		if err != nil {
			return err
		}
		// secret can be use to encryption
		// encryptedAppJSON, err := utils.EncryptData(appJSON, s.secret)
		// if err != nil {
		// 	return err
		// }
		return txn.Set([]byte(app_prefix+app.AppDID), appJSON)
	})
}

// GetApp retrieves an Application instance from the database
func (s *Store) GetApp(appID string) (*models.ApplicationResponse, error) {
	var app models.ApplicationResponse
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(app_prefix + appID))
		if err != nil {
			return err
		}

		return item.Value(func(val []byte) error {
			// // Decrypt data
			// decryptedVal, err := utils.DecryptData(val, s.secret)
			// if err != nil {
			// 	return err
			// }
			return json.Unmarshal(val, &app)
		})
	})
	if err != nil {
		return nil, err
	}
	return &app, nil
}

// SetAuthProvider stores an AuthProvider instance in the database
func (s *Store) SetAuthProvider(auth models.AuthProvider) error {
	return s.db.Update(func(txn *badger.Txn) error {
		authJSON, err := json.Marshal(auth)
		if err != nil {
			return err
		}
		return txn.Set([]byte(auth_prefix+auth.AppDID), authJSON)
	})
}

// GetAuthProvider retrieves an AuthProvider instance from the database
func (s *Store) GetAuthProvider(appID string) (*models.AuthProvider, error) {
	var auth models.AuthProvider
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(auth_prefix + appID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &auth)
		})
	})
	if err != nil {
		return nil, err
	}
	return &auth, nil
}

// SetPolicy stores a PolicySchemaResponse instance in the database
func (s *Store) SetPolicy(policy models.PolicySchemaResponse) error {
	return s.db.Update(func(txn *badger.Txn) error {
		policyJSON, err := json.Marshal(policy)
		if err != nil {
			return err
		}
		return txn.Set([]byte(policy_prefix+policy.ID), policyJSON)
	})
}

// GetPolicy retrieves a PolicySchemaResponse instance from the database
func (s *Store) GetPolicy(policyID string) (*models.PolicySchemaResponse, error) {
	var policy models.PolicySchemaResponse
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(policy_prefix + policyID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &policy)
		})
	})
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// SetIssuedPolicy stores an ApplicationPolicyResponse instance in the database
func (s *Store) SetIssuedPolicy(policy models.ApplicationPolicyResponse) error {
	return s.db.Update(func(txn *badger.Txn) error {
		policyJSON, err := json.Marshal(policy)
		if err != nil {
			return err
		}
		return txn.Set([]byte(issued_policy_prefix+policy.ApplicationDID), policyJSON)
	})
}

// GetIssuedPolicy retrieves an ApplicationPolicyResponse instance from the database
func (s *Store) GetIssuedPolicy(appDID string) (*models.ApplicationPolicyResponse, error) {
	var policy models.ApplicationPolicyResponse
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(issued_policy_prefix + appDID))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &policy)
		})
	})
	if err != nil {
		return nil, err
	}
	return &policy, nil
}

// GetAllApps retrieves all Application instances from the database
func (s *Store) GetAllApps() ([]models.ApplicationResponse, error) {
	var apps []models.ApplicationResponse
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(app_prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			if !bytes.HasPrefix(key, []byte(app_prefix)) {
				continue
			}
			err := item.Value(func(val []byte) error {
				var app models.ApplicationResponse
				if err := json.Unmarshal(val, &app); err != nil {
					return err
				}
				apps = append(apps, app)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return apps, nil
}

// GetAllPolicies retrieves all PolicySchemaResponse instances from the database
func (s *Store) GetAllPolicies() ([]models.PolicySchemaResponse, error) {
	var policies []models.PolicySchemaResponse
	err := s.db.View(func(txn *badger.Txn) error {
		opts := badger.DefaultIteratorOptions
		opts.Prefix = []byte(policy_prefix)
		it := txn.NewIterator(opts)
		defer it.Close()

		for it.Rewind(); it.Valid(); it.Next() {
			item := it.Item()
			key := item.Key()
			if !bytes.HasPrefix(key, []byte(policy_prefix)) {
				continue
			}
			err := item.Value(func(val []byte) error {
				var policy models.PolicySchemaResponse
				if err := json.Unmarshal(val, &policy); err != nil {
					return err
				}
				policies = append(policies, policy)
				return nil
			})
			if err != nil {
				return err
			}
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return policies, nil
}

// SetProviderSchema sets provider schema details.
func (s *Store) SetProviderSchema(prov models.ProviderSchema) error {
	return s.db.Update(func(txn *badger.Txn) error {
		appJSON, err := json.Marshal(prov)
		if err != nil {
			return err
		}
		return txn.Set([]byte(provider_schema_prefix+prov.ProviderName), appJSON)
	})
}

// GetProviderSchema get provider schema details.
func (s *Store) GetProviderSchema(provider string) (*models.ProviderSchema, error) {
	var prov models.ProviderSchema
	err := s.db.View(func(txn *badger.Txn) error {
		item, err := txn.Get([]byte(provider_schema_prefix + provider))
		if err != nil {
			return err
		}
		return item.Value(func(val []byte) error {
			return json.Unmarshal(val, &prov)
		})
	})
	if err != nil {
		return nil, err
	}
	return &prov, nil
}
