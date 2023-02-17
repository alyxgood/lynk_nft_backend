# lynk_nft_backend

### Config
1. RPC: https://github.com/lynkdevelopment/lynk_nft_backend/blob/main/consts/const.go#L13
2. LYNKNFT Address: https://github.com/lynkdevelopment/lynk_nft_backend/blob/main/consts/const.go#L14
3. Metadata.image: https://github.com/lynkdevelopment/lynk_nft_backend/blob/main/consts/const.go#L15
4. Server port: https://github.com/lynkdevelopment/lynk_nft_backend/blob/main/consts/const.go#L16
5. Metadata cache: https://github.com/lynkdevelopment/lynk_nft_backend/blob/main/main.go#L58

### Deploy
```shell
cd lynk_nft_backend
./utils/build_binaries.sh
nohup ./bin/alyx_nft_backend &
```