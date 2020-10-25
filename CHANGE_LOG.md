# Change Log

## [v1.5](https://github.com/thewizardplusplus/go-hashmap/tree/v1.5) (2020-10-25)

- support stopping of iteration via a context:
  - for a hash map;
  - for a synchronized hash map;
  - for a concurrent hash map.

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration:
        - via a handling result;
        - via a context;
      - support randomizing of iteration order;
    - setting of an item by a key;
    - deleting of an item by a key;
  - support options:
    - initial capacity;
    - maximal load factor;
    - grow factor;
- implementation of a synchronized hash map:
  - use the interface of an universal storage as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration:
        - via a handling result;
        - via a context;
      - support randomizing of iteration order;
    - setting of an item by a key;
    - deleting of an item by a key;
  - support options:
    - inner map;
- implementation of a concurrent hash map:
  - use data sharding for concurrent access;
  - use the interface of an universal storage as one shard;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration:
        - via a handling result;
        - via a context;
      - support randomizing of iteration order:
        - over items and their keys;
        - over shards;
    - setting of an item by a key;
    - deleting of an item by a key;
  - support options:
    - concurrency level;
    - shard factory.

## [v1.4.1](https://github.com/thewizardplusplus/go-hashmap/tree/v1.4.1) (2020-10-16)

- add explanatory comments for implementations:
  - of a hash map;
  - of a synchronized hash map;
  - of a concurrent hash map;
- describe for releases:
  - features;
  - change log.

## [v1.4](https://github.com/thewizardplusplus/go-hashmap/tree/v1.4) (2020-10-12)

- add the interface of an universal storage:
  - implementation of a synchronized hash map:
    - use the interface of an universal storage as an inner map;
  - implementation of a concurrent hash map:
    - use the interface of an universal storage as one shard;
- support options:
  - for a hash map:
    - initial capacity;
    - maximal load factor;
    - grow factor;
  - for a synchronized hash map:
    - inner map;
  - for a concurrent hash map:
    - concurrency level;
    - shard factory;
- misc.:
  - pass the `hashmap.Key` interface as a parameter to the `hashmap.Key.Equals()` method;
  - improve the Travis CI configuration.

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
      - support randomizing of iteration order;
    - setting of an item by a key;
    - deleting of an item by a key;
  - support options:
    - initial capacity;
    - maximal load factor;
    - grow factor;
- implementation of a synchronized hash map:
  - use the interface of an universal storage as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
      - support randomizing of iteration order;
    - setting of an item by a key;
    - deleting of an item by a key;
  - support options:
    - inner map;
- implementation of a concurrent hash map:
  - use data sharding for concurrent access;
  - use the interface of an universal storage as one shard;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
      - support randomizing of iteration order:
        - over items and their keys;
        - over shards;
    - setting of an item by a key;
    - deleting of an item by a key;
  - support options:
    - concurrency level;
    - shard factory.

## [v1.3](https://github.com/thewizardplusplus/go-hashmap/tree/v1.3) (2019-06-28)

- support randomizing of iteration order:
  - for a hash map:
    - over items and their keys;
  - for a synchronized hash map:
    - over items and their keys;
  - for a concurrent hash map:
    - over items and their keys;
    - over shards.

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
      - support randomizing of iteration order;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a synchronized hash map:
  - use the implementation described above as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
      - support randomizing of iteration order;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a concurrent hash map:
  - use data sharding for concurrent access;
  - use the synchronized implementation described above as one shard;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
      - support randomizing of iteration order:
        - over items and their keys;
        - over shards;
    - setting of an item by a key;
    - deleting of an item by a key.

## [v1.3-alpha.2](https://github.com/thewizardplusplus/go-hashmap/tree/v1.3-alpha.2) (2019-06-26)

- support stopping of iteration over items and their keys:
  - for a hash map;
  - for a synchronized hash map;
  - for a concurrent hash map.

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a synchronized hash map:
  - use the implementation described above as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a concurrent hash map:
  - use data sharding for concurrent access;
  - use the synchronized implementation described above as one shard;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys:
      - support stopping of iteration;
    - setting of an item by a key;
    - deleting of an item by a key.

## [v1.3-alpha.1](https://github.com/thewizardplusplus/go-hashmap/tree/v1.3-alpha.1) (2019-05-11)

- support iteration over items and their keys:
  - for a hash map;
  - for a synchronized hash map;
  - for a concurrent hash map.

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a synchronized hash map:
  - use the implementation described above as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a concurrent hash map:
  - use data sharding for concurrent access;
  - use the synchronized implementation described above as one shard;
  - support operations:
    - getting of an item by a key;
    - iteration over items and their keys;
    - setting of an item by a key;
    - deleting of an item by a key.

## [v1.2](https://github.com/thewizardplusplus/go-hashmap/tree/v1.2) (2019-05-04)

- implementation of a concurrent hash map:
  - use data sharding for concurrent access;
  - use the synchronized implementation described above as one shard;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key;
- remove the success flag from deleting methods.

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a synchronized hash map:
  - use the implementation described above as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key;
- implementation of a concurrent hash map:
  - use data sharding for concurrent access;
  - use the synchronized implementation described above as one shard;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key.

## [v1.1](https://github.com/thewizardplusplus/go-hashmap/tree/v1.1) (2019-05-02)

- implementation of a synchronized hash map:
  - use the implementation described above as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key:
      - support a success flag.

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key:
      - support a success flag;
- implementation of a synchronized hash map:
  - use the implementation described above as an inner map;
  - use a mutex lock to access the inner map;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key:
      - support a success flag.

## [v1.0](https://github.com/thewizardplusplus/go-hashmap/tree/v1.0) (2019-05-01)

### Features

- implementation of a hash map:
  - use the open addressing strategy for collision resolution;
  - use the key interface for supporting custom types;
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key:
      - support a success flag.
