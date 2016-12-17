package xisdb

import "github.com/alexsward/xisdb/ql"

// QueryEngine is the processor of xisql statements
type QueryEngine struct {
}

// QueryEngineContext is the context of what the engine should perform the query in
type QueryEngineContext struct {
	DB      *DB
	Results chan Item
	// errors  chan<- error
}

// Execute will perform all of the statements in the context of the QueryEngineContext
func (qe *QueryEngine) Execute(statements []ql.Statement, ctx *QueryEngineContext) error {
	if ctx.DB == nil {
		return ErrNoDatabase
	}

	go func() error {
		defer close(ctx.Results)
		for _, statement := range statements {
			if err := statement.Validate(); err != nil {
				return err
			}

			switch statement.(type) {
			case *ql.GetStatement:
				s := statement.(*ql.GetStatement)
				for _, key := range s.Keys() {
					item, err := ctx.DB.Get(key)
					if err != nil {
						return err
					}
					ctx.Results <- Item{key, item, nil}
				}
				return nil
			case *ql.SetStatement:
				s := statement.(*ql.SetStatement)
				return ctx.DB.ReadWrite(func(tx *Tx) error {
					for key, value := range s.Pairs() {
						err := tx.Set(key, value, nil)
						if err != nil {
							return err
						}
						ctx.Results <- Item{key, value, nil}
					}
					return nil
				})
			case *ql.DelStatement:
				s := statement.(*ql.DelStatement)
				return ctx.DB.ReadWrite(func(tx *Tx) error {
					for _, key := range s.Keys() {
						_, err := tx.Delete(key)
						if err != nil {
							return err
						}
						ctx.Results <- Item{key, "", nil}
					}
					return nil
				})
			}
		}
		return nil
	}()
	return nil
}
