# Orm

The Orm module is designed to help you easily set up and manage your gorm instance using configuration files. This module is highly modular and can be integrated seamlessly with other tools and services, such as Otel, NewRelic, etc. It is compatible with any SQL database. The module provides key components including orm.GormProvider, *gorm.DB, and tx.Transactioner, each serving a specific purpose to streamline database operations in your application.

## GormProvider

GormProvider is a utility that simplifies the process of retrieving a *gorm.DB instance. It is the preferred method for accessing the gorm.DB instance, rather than injecting it directly via a constructor.


The Get method in GormProvider checks the context for an existing transaction. If a transaction is found, it returns that specific instance. Otherwise, it provides a gorm.DB instance with the given context
This allows for consistent and context-aware access to your database, ensuring that operations respect transactional boundaries when necessary.
```golang
func (p GormProvider) Get(ctx context.Context) *gorm.DB {
	if found := tx.From(ctx); found != nil {
		return found
	}
	return p.db.WithContext(ctx)
}
```

## Transactioner
Transactioner is a utility struct designed to facilitate transaction management within gorm by utilizing context. This is particularly useful in services where multiple operations need to be executed as part of a single transaction.


### Example Usage

The following is an example of a service that uses Transactioner to manage a series of database operations within a single transaction:

```golang
type Service struct {
    txr tx.Transactioner
    repo Repository
}

func NewService(txr tx.Transactioner, repo Repositroy) *Service {
    return &service{txr: txr, repo: repo}
}

func (s *service) Save(parentctx context.Context) error {
    return s.txr.Transaction(parentctx, func(ctx context.Context) error {
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
In this example, the Save method encapsulates multiple operations (SaveOne and SaveTwo) within a single transaction. If any operation fails, the entire transaction is rolled back, maintaining data integrity.


## Drivers
Drivers in this module provide an abstraction layer for different gorm or sql/db dialects, allowing the Orm module to work with various SQL databases. Each driver can have its own unique configuration and options, making it flexible to adapt to different database setups.


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


## Repository

The Repository interface is an abstraction built on top of gorm, designed to minimize repetitive code for common database operations. This interface provides a standardized set of methods for interacting with the database, such as finding, saving, updating, and deleting records.

```go

type Repository[Id IdType, T any] interface {
	FindById(ctx context.Context, id Id) (*T, error)
	FindOne(ctx context.Context, q query.Query) (*T, error)
	FindAll(ctx context.Context, q query.Query) ([]*T, error)
	Update(ctx context.Context, q query.Query, update any) error
	Save(ctx context.Context, t *T) (*T, error)
	SaveAll(ctx context.Context, ts []*T) error
	DeleteById(ctx context.Context, id Id) error
	Delete(ctx context.Context, q query.Query) error
	ListPageable(ctx context.Context, q query.Query, req PageableRequest) (*PageableResponse[*T], error)
	ParseQuery(ctx context.Context, q query.Query) *gorm.DB
	Context(ctx context.Context) *gorm.DB
}

```

### Usage
```golang
// Model
type Foo struct {
    Id int `gorm:"primaryKey"`
    Bar string
}

func NewRepository(gp orm.GormProvider) repository.Repository[int,Foo] {
    return repository.New[int, Foo](gp)
}

type Service struct {
    fooRepo repository.Repository[int,Foo]
}

func (s *service) GetFoo(ctx context.Context, id int) (*Foo, error) {
    return s.fooRepo.FindById(ctx, id)
}

```

### With Composition
Composition allows you to extend the base repository with custom methods, giving you the flexibility to add specific queries or operations while still leveraging the core repository functionality.

```golang
type FooRepository struct {
    repository.Repository[int, Foo]
}

func NewRepository(gp gorm.Provider) *FooRepository {
    return &FooRepository{
        Repository: repository.New[int, Foo]
    }
}

// Custom repo function
func (repo *FooRepository) AdvenceQuery(ctx context.Context, id int) (*Foo, error) {
    // you can acces all base class functionalities here
    return repo.FindById(ctx, id)
}

func (repo *FooRepository) DirectGormQuery(ctx context.Context) error {
    // you can get gorm instance here
    gormdb := repo.Context(ctx)

    // gorm operations...

    return nil
}

```

In this case, the FooRepository struct extends the base repository and adds custom methods like AdvenceQuery and DirectGormQuery, providing additional functionality while maintaining access to the core repository methods.