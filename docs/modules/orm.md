# Orm

The Orm module allows you to build your `gorm` instance using config files. and easily pluggable by other modules such as Otel, NewRelic, etc. The module runs with any SQL DB. This module provides `orm.GormProvider`, `*gorm.DB`, and `tx.Transactioner`

## GormProvider
GormProvider is a basic `*gorm.DB` instance getter. İt is the recommended way to access Gorm.DB instance instead of directly injected from a constructor

```golang
func (p GormProvider) Get(ctx context.Context) *gorm.DB {
	if found := tx.From(ctx); found != nil {
		return found
	}
	return p.db.WithContext(ctx)
}
```

## Transactioner
Transactioner is a utility struct for handling `gorm` transactions with context
```golang
type Service struct {
    txr tx.Transactioner
    repo Repository
}

func NewService(txr tx.Transactioner, repo Repositroy) *Service {
    return &service{txr: txr, repo: repo}
}

func (s *service) Save(ctx context.Context) error {
    return s.txr.Transaction(func(ctx context.Context) error {
        if err := repo.SaveOne(ctx); err != nil {
            return err
        }

        if err := repo.SaveTwo(ctx); err != nil {
            return err
        }

        return nil
    })
}
```


## Drivers
Driver is an abstraction of the ORM module to use any `gorm` or `sql/db` dialect. Each driver can have its unique configuration and options. Drivers must be passed on module constructor.


### Postgres
It is a `gorm` Postgres driver abstraction for Chaki

#### Usage
```golang
func main() {
    app := chaki.New()

    app.Use(
        orm.Module(postgresdriver.New()),
    )

    app.Provide(
        newRepository,
    )
}


type Repository struct {
    gp orm.GormProvider
}

func newRepository(gp orm.GormProvider) *Repository {
    return &repository{
        gp: gp,
    }
}

func (repo *repository) someOperation(ctx context.Context) {
    db := repo.gp.Get(ctx)
    //  you can use gorm instance
}

```

#### Config
```yaml
gorm:
  host: "string" #required
  port: "number" #required
  username: "string" #required
  password: "string" #required

```


