# 3scale Operator

## 3scale API Management installation functionality

The following Custom Resources are provided:

`APIManager`

One APIManager custom resource per project is allowed.

This resource is the resource used to deploy a 3scale API Management solution.

### APIManager

| **Field** | **json/yaml field**| **Type** | **Required** | **Description** |
| --- | --- | --- | --- | --- |
| Spec | `spec` | [APIManagerSpec](#APIManagerSpec) | The specfication for APIManager custom resource | Yes |
| Status | `status` | [APIManagerStatus](#APIManagerStatus) | The status for the custom resource | No |

#### APIManagerSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| ProductVersion | `productVersion` | ProductVersion | Yes | N/A | 3Scale API Management solution release version. Currently the values "2.5" and "upstream" are supported. **Warning**, "upstream" uses 3scale unstable nightly images and it is intended only for development purposes and to be used in non-productive environments |
| WildcardDomain | `wildcardDomain` | string | Yes | N/A | Root domain for the wildcard routes. Eg. example.com will generate 3scale-admin.example.com. |
| AppLabel | `appLabel` | string | No | `3scale-api-management` | The value of the `app` label that will be applied to the API management solution
| TenantName | `tenantName` | string | No | `3scale` | Tenant name under the root that Admin UI will be available with -admin suffix.
| WildcardPolicy | `wildcardPolicy` | string | No | `None` | Use `Subdomain` to create a wildcard route for apicast wildcard router. If `Subdomain` is used, wildcard routes at the OpenShift router level need to be enabled. You can do so by executing `oc set env dc/router ROUTER_ALLOW_WILDCARD_ROUTES=true -n default` |
| ImageStreamTagImportInsecure | `imageStreamTagImportInsecure` | bool | No | `false` | Set to true if the server may bypass certificate verification or connect directly over HTTP during image import |
| ResourceRequirementsEnabled | `resourceRequirementsEnabled` | bool | No | `true` | When true, 3Scale API management solution is deployed with the optimal resource requirements and limits. Setting this to false removes those resource requirements. ***Warning*** Only set it to false for development and evaluation environments |
| ApicastSpec | `apicast` | \*ApicastSpec | No | See [ApicastSpec](#ApicastSpec) | Spec of the Apicast part |
| BackendSpec | `backend` | \*BackendSpec | No | See [BackendSpec](#BackendSpec) reference | Spec of the Backend part |
| SystemSpec  | `system`  | \*SystemSpec  | No | See [SystemSpec](#SystemSpec) reference | Spec of the System part |
| WildcardRouterSpec | `wildcardRouter` | \*WildcardRouterSpec | No | See [WildcardRouterSpec](#WildcardRouterSpec) reference | Spec of the WildcardRouter part |
| ZyncSpec    | `zync`    | \*ZyncSpec    | No | See [ZyncSpec](#ZyncSpec) reference | Spec of the Zync part    |
| HighAvailabilitySpec | `highAvailability` | \*HighAvailabilitySpec | No | See [HighAvailabilitySpec](#HighAvailabilitySpec) reference | Spec of the HighAvailability part |

#### ApicastSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| ApicastManagementAPI | `managementAPI` | string | No | `status` | Scope of the APIcast Management API. Can be disabled, status or debug. At least status required for health checks |
| OpenSSLVerify | `openSSLVerify` | bool | No | `false` | Turn on/off the OpenSSL peer verification when downloading the configuration |
| IncludeResponseCodes  | `responseCodes` | bool | No | `true` | Enable logging response codes in APIcast |
| RegistryURL | `registryURL` | string | No | `http://apicast-staging:8090/policies` | The URL to point to APIcast policies registry management |
| Image | `image` | string | No | nil | Used to overwrite the desired container image for Apicast |

#### BackendSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| Image | `image` | string | No | nil | Used to overwrite the desired container image for Backend |
| RedisImage | `redisImage` | string | No | nil | Used to overwrite the desired Redis image for the Redis used by backend |

#### SystemSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| Image | `image` | string | No | nil | Used to overwrite the desired container image for System |
| RedisImage | `redisImage` | string | No | nil | Used to overwrite the desired Redis image for the Redis used by System |
| MemcachedImage | `memcachedImage` | string | No | nil | Used to overwrite the desired Memcached image for the Memcached used by System |
| FileStorageSpec | `fileStorage` | \*SystemFileStorageSpec | No | See [FileStorageSpec](#FileStorageSpec) specification | Spec of the System's File Storage part |
| DatabaseSpec | `database` | \*SystemDatabaseSpec | No | See [SystemDatabaseSpec](#SystemDatabaseSpec) specification | Spec of the System's Database part |

#### FileStorageSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| PVC | `persistentVolumeClaim` | \*SystemPVCSpec | No | nil | Used to use a PersistentVolumeClaim as the System's file storage. See [SystemPVCSpec](#SystemPVCSpec) |
| S3  | `amazonSimpleStorageService` | \*SystemS3Spec | No | nil | Used to use S3 as the System's file storage. See [SystemS3Spec](#SystemS3Spec) |

Only one of the fields can be chosen. If no field is specified then PVC is used.

#### SystemPVCSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| StorageClassName | `storageClassName` | string | No | nil | The Storage Class to be used by the PVC |

#### SystemS3Spec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| AWSBucket | `awsBucket` | string | Yes | N/A | AWS Bucket name of the S3 bucket to be used as System's FileStorage for assets |
| AWSRegion | `awsRegion` | string | Yes | N/A | AWS Region of the S3 bucket to be used as Sytem's FileStorage for assets |
| AWSCredentials | `awsCredentialsSecret` | [corev1.LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#localobjectreference-v1-core) | Yes | N/A | Local object reference to the secret to be used where the AWS credentials are stored. See [LocalObjectReference](https://kubernetes.io/docs/reference/generated/kubernetes-api/v1.11/#localobjectreference-v1-core) on how to specify the local object reference to the secret |
| FileUploadStorage | `fileUploadStorage` | string | Yes | N/A | Define Assets Storage name |

The secret name specified in the `awsCredentialsSecret` field must be
pre-created by the user before creating the APIManager custom resource.
Otherwise the operator will complain about it. See the
[fileStorage S3 credentials secret](#fileStorage-S3-credentials-secret)
specification to see what fields the secret should have and the values
that should be set on it.

#### DatabaseSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| MySQLSpec | `mysql`| \*SystemMySQLSpec | No | nil | See [MySQLSpec](#MySQLSpec) specification |

#### MySQLSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| Image | `image` | string | No | nil | Used to overwrite the desired container image for System's MySQL database |

#### WildcardRouterSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| Image | `image` | string | No | nil | Used to overwrite the desired container image for WildcardRouter |

#### ZyncSpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| Image | `image` | string | No | nil | Used to overwrite the desired container image for Zync |
| PostgreSQLImage | `postgreSQLImage` | string | No | nil | Used to overwrite the desired PostgreSQL image for the PostgreSQL used by Zync |

#### HighAvailabilitySpec

| **Field** | **json/yaml field**| **Type** | **Required** | **Default value** | **Description** |
| --- | --- | --- | --- | --- | --- |
| Enabled | `enabled` | bool | No | `false` | Enable to use external system database, backend redis, system redis and apicast redis databases|

When HighAvailability is enabled the following secrets have to been
pre-created by the user:

* [backend-redis](#backend-redis) with the `REDIS_STORAGE_URL` and
  `REDIS_QUEUES_URL` fields with values pointing to the desired external
  databases. The databases should be configured
  in high-availability mode
* [system-database](#system-database) with the `URL` field with the value
  pointing to the desired external database. The database should be configured
  in high-availability mode
* [system-redis](#system-redis) with the `URL` and `MESSAGE_BUS_URL` fields
  with the value pointing to the desired external databases. The databases
  should be configured in high-availability mode

#### APIManagerStatus

Used by the Operator/Kubernetes to control the state of the APIManager.
an `APIManager` status field should never be modified by the user.

| **Field** | **json/yaml field**| **Type** | **Info** |
| --- | --- | --- | --- |
| No fields for the moment | | | |

### APIManager Secrets

Additionally, if desired, several sensitive APIManager configuration options
can be controlled by pre-creating some Kubernetes secrets before deploying the
APIManager Custom Resource.

The available configurable secrets are:

#### backend-internal-api

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| username | Backend internal API username. Backend internal API is used by System | `3scale_api_user` |
| password | Backend internal API password. Backend internal API is used by System | Autogenerated value |

#### backend-listener

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| service_endpoint | Backend listener service endpoint. Used by System | `http://backend-listener:3000` |
| route_endpoint | Backend listener route endpoint. Used by System | `https://backend-<tenantName>.<wildcardDomain>` |

#### backend-redis

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| REDIS_STORAGE_URL | Backend's redis storage database URL | `redis://backend-redis:6379/0` |
| REDIS_STORAGE_SENTINEL_ROLE | Backend's redis storage sentinel role name. Used only when Redis sentinel is configured in the Redis database being used | `""` |
| REDIS_STORAGE_SENTINEL_HOSTS | Backend's redis storage sentinel hosts name. Used only when Redis sentinel is configured in the Redis database being used | `""` |
| REDIS_QUEUES_URL | Backend's redis queues database URL  | `redis://backend-redis:6379/1` |
| REDIS_QUEUES_SENTINEL_ROLE | Backend's redis queues sentinel role name. Used only when Redis sentinel is configured in the Redis database being used | `""` |
| REDIS_QUEUES_SENTINEL_HOSTS | Backend's redis queues sentinel hosts name. Used only when Redis sentinel is configured in the Redis database being used | `""` |

#### system-app

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| SECRET_KEY_BASE | System application secret key base | Autogenerated value |

#### system-database

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| URL | URL of the Porta database. The format of the URL must be: `mysql2://root:<RootPassword>@<DatabaseHost>/<DatabaseName>` | `mysql2://root:<AutogeneratedValue>@system-mysql/<AutogeneratedValue>` where '<>' fields should be replaced by the desired values |
| DB_USER | Non-administrative database username | `mysql` |
| DB_PASSWORD | Password of the non-administrative database user | Autogenerated value |

#### system-events-hook

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| URL | TODO | `http://system-master:3000/master/events/import` |
| PASSWORD | Shared secret to import events from backend to system | Autogenerated value |

#### system-master-apicast

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| ACCESS_TOKEN | Read only access token that APIcast uses to download its configuration from System | Autogenerated value |
| BASE_URL | URL of the 3scale portal admin endpoint with authentication part | `http://<ACCESS_TOKEN>@system-master:3000` |
| PROXY_CONFIGS_ENDPOINT | URL of the available configs for all System's services | `http://<ACCESS_TOKEN>@system-master:3000/master/api/proxy/configs` |

#### system-memcache

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| SERVERS | System's Memcached URL | `system-memcache:11211` |

#### system-recaptcha

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| PUBLIC_KEY | reCAPTCHA site key (used in spam protection) for System| `""` |
| SECRET_KEY | reCAPTCHA secret key (used in spam protection) for System| `""` |

#### system-redis

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| URL | System's Redis database URL | `redis://system-redis:6379/1` |
| MESSAGE_BUS_URL | System's Message Bus Redis database URL | `redis://system-redis:6379/8` |
| NAMESPACE | Define the namespace to be used by System's Redis Database. The empty value means not namespaced | `""` |
| MESSAGE_BUS_NAMESPACE | Define the namespace to be used by System's Message Bus Redis Database. The empty value means not namespaced | `""` |

#### system-seed

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| MASTER_USER | System's master username | `master` |
| MASTER_PASSWORD | System's master password | Autogenerated value |
| MASTER_ACCESS_TOKEN | System's master access token | Autogenerated value |
| MASTER_DOMAIN | System's master domain name | `master` |
| ADMIN_USER | System's admin user of the tenant created by default | `admin` |
| ADMIN_PASSWORD | System's admin password of the tenant created by default | Autogenerated value |
| ADMIN_ACCESS_TOKEN | System's admin access token of the tenant created by default | Autogenerated value |
| TENANT_NAME | Tenant name under the root that Admin UI will be available with -admin suffix | `<tenantName>` |

#### zync

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| DATABASE_URL | PostgreSQL database used by Zync. | `postgresql://zync:<ZYNC_DATABASE_PASSWORD>@zync-database:5432/zync_production` |
| SECRET_KEY_BASE | Zync's application key generator to encrypt communications | Autogenerated value |
| ZYNC_AUTHENTICATION_TOKEN | Authentication token used to authenticate System when calling Zync | Autogenerated value |
| ZYNC_DATABASE_PASSWORD | Database password associated to the 'zync' user (non-admin user) | Autogenerated value |

#### fileStorage-S3-credentials-secret

The name of this secret can be any name as long as does not collide with other
existing secret names.

| **Field** | **Description** | **Default value** |
| --- | --- | --- |
| AWS_ACCESS_KEY_ID | AWS Access Key ID to use in S3 Storage for System's file storage | N/A |
| AWS_SECRET_ACCESS_KEY | AWS Access Key Secret to use in S3 Storage for System's file storage | N/A |
