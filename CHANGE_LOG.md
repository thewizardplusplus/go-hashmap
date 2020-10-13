# Change Log

## [v1.4](https://github.com/thewizardplusplus/go-hashmap/tree/v1.4) (2020-10-12)

## [v1.3](https://github.com/thewizardplusplus/go-hashmap/tree/v1.3) (2019-06-28)

## [v1.3-alpha.2](https://github.com/thewizardplusplus/go-hashmap/tree/v1.3-alpha.2) (2019-06-26)

## [v1.3-alpha.1](https://github.com/thewizardplusplus/go-hashmap/tree/v1.3-alpha.1) (2019-05-11)

## [v1.2](https://github.com/thewizardplusplus/go-hashmap/tree/v1.2) (2019-05-04)

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
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key:
      - support a success flag;
  - use the key interface for supporting custom types;
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
  - support operations:
    - getting of an item by a key;
    - setting of an item by a key;
    - deleting of an item by a key:
      - support a success flag;
  - use the key interface for supporting custom types.
