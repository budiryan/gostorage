# GoStorage
Similar to <a href="https://github.com/febytanzil/gobroker">gobroker</a>, but for cloud storage.

## Install
### Standard
`go get github.com/budiryan/gostorage/storage`

### Dep
`dep ensure -add github.com/budiryan/gostorage/storage`

## Currently Supported Cloud Storage Providers
1. <a href="https://cloud.google.com/storage/">Google Cloud Storage</a>

## Currently Supported Operations
1. Read
2. Write
3. Check Object / Blob's Existence
4. Get Signed URL
5. List Object Inside a Bucket

## Usage Examples
Please see <a href="./example_test.go">example_test.go</a> for example usages

## Mocking
Mock functionality is provided for your unit-tests. See <a href="./storage/mock_storage/mock_storage.go">here</a>
