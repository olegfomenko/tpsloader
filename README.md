---
    title: Load testing for max tps
    author: olegfomenko 
    date: 22 January 2021 
---

# TPSLoader - load testing tool

TPSLoader it is a load testing tool for stellar + horizon network that will produce CreateAccount and Payment transaction load

For correct work and simplified configuration service uses pool of 
"creators" - accounts that will pay their lumen for CreateAccount Transaction and 
"payers" - pair of accounts that will pay each other small count of lumen

## Config
``` yaml
admin: SDR7XY33FYTDJTRF2CAXU5VGQWIQU4YOGZYPMYZ7ZAZTDGINQYMRJWZC

passphrase: Stellar Load Test Network

horizon: http://host.docker.internal:8000/

creators: ["SAN4NLOGBGWHNJADLOSWT4GBYKMP7ZMGE74NLMHHRFXL6PB45HVWLX3W", "SAMEBGASYZEAJZLTKMWFUE47ZMPNIP34WNHAEYN5TVF7DZQW3242XHKM", "SC4JZ3ML5KCJFY43L5YONXWUETHF6WNMCGQGM4VTP7G3SH6HEWCWWBQJ"]

payers:
  SCQJXD46Y2SSTYIUO6XMFUXL6UAQAACDLU75QPAIH7XRWLJGQKTXG2KT: SB2RJMYSGQVYJ42TKJZK6HNF3O6ULJV5Y4I2TRGQO5SE2MRKESBI77YG
  SBQLTUKFCNJPVQOKYX5LAGIXIOOHT62V5I3MKFNA3MA4UXPWZTTJKYCD: SCJNF47C6WZ52NOMBNZ6SGHGGXWDCRIUXXOUWNWBKZLLBHHO47SO73JC

amount: 10

duration: 600000
```

- admin -> stellar network admin seed (for creating accounts)
- passphrase -> core passphrase
- horizon -> horizon microservice URL
- creators -> array of "creators" (not required for auto-run startup)
- payers -> map of "payers" (not required for auto-run startup)
- amount -> count of lumen that payers will pay each other
- duration -> load testing duration (does not include creators and payment account creation)

## Commands

### get-acc

gen-acc -n N -a A -> command for generation accounts (not a load testing) for future using in configuration

-n N -> number of accounts to create

-a A -> amount for pay new-created accounts

Example:
```commandline
./tpsloader gen-acc -n 5 -a 200
``` 

### run
run -> command for start load testing (creators and payers config required!)

### auto-run

auto-run -n N -a1 A1 -p P -a2 A2 -> command for run load testing with creators and payer initialization before. 
You don't need to configure in yourself, service will create it from admin account and run load testing based on created accounts

Example:
```commandline
./tpsloader auto-run -n 500 -a1 100000 -p 500 -a2 200
``` 
