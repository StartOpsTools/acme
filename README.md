ACME
====

通过了解 ACME 标准, 了解自动办法证书过程. 简单重新制作一个 ACME 客户端轮子.

未成功对接

[知乎参考链接](https://zhuanlan.zhihu.com/p/75032510)

[ACME rfc](https://datatracker.ietf.org/doc/html/rfc8555)


request Header：
    Accept-Language

BASE64URL




JWS 保护标头必须包含以下字段：

      *“alg”（算法）

         + 此字段不得包含“无”或消息
            验证码 (MAC) 算法（例如，其中的一种
            算法注册表描述提到 MAC/HMAC）。

      * “nonce”（在第 6.5 节中定义）

      * “url”（在6.4 节中定义）



“jwk”和“kid”字段是互斥的。服务器必须拒绝包含两者的请求。


对于 newAccount 请求，以及通过身份验证的 revokeCert 请求证书密钥，必须有一个“jwk”字段。该字段必须包含与用于签名的私钥对应的公钥 JWS。

   

一旦超过速率限制，  服务器必须响应类型错误 “urn:ietf:params:acme:error:rateLimited”。此外，服务  应该发送一个 Retry-After 头域 [ RFC7231 ] 指示何时   当前请求可能会再次成功。



ACME v2 (RFC 8555)
[生产环境] https://acme-v02.api.letsencrypt.org/directory
[测试环境] https://acme-staging-v02.api.letsencrypt.org/directory




ACME 支持：

   o 帐户创建

   o 订购证书

   o 标识符授权

   o 证书颁发

   o 证书吊销
   
   

https://github.com/go-acme/lego



在 ACME 中，可以创建一个帐户并将其用于所有授权和颁发，或者为每个客户创建一个帐户。


在使用 dns-01 进行验证时，请务必清理旧的 TXT 记录，以确保对 Let’s Encrypt 的查询的回复不会太大。

我们建议在有效期剩余三分之一时自动续期证书。由于 Let’s Encrypt 当前颁发有效期为 90 天的证书，这意味着在到期前 30 天进行续期。

请特别确保这些客户端定期收到更新。将来，这些事情可能会发生变化：
    
    我们用于签发证书的根证书和中间证书
    对证书签名时使用的哈希算法
    我们愿意用于对终端实体证书进行签名的密钥类型和密钥强度
    ACME 协议

中间证书:
    
    你不应该硬编码中间证书，而应该使用 ACME 协议中的 Link: rel="up" 标头，因为中间证书很可能会改变。

服务条款（ToS）:
    
    请避免硬编码 ToS 链接，而是使用 Link：rel =“terms-of-service” 标头确定要使用的 ToS 链接。





## Step

1. /directory

GET /directory

```json
{
  "4bmRla_FAhE": "https://community.letsencrypt.org/t/adding-random-entries-to-the-directory/33417",
  "keyChange": "https://acme-v02.api.letsencrypt.org/acme/key-change",
  "meta": {
    "caaIdentities": [
      "letsencrypt.org"
    ],
    "termsOfService": "https://letsencrypt.org/documents/LE-SA-v1.2-November-15-2017.pdf",
    "website": "https://letsencrypt.org"
  },
  "newAccount": "https://acme-v02.api.letsencrypt.org/acme/new-acct",
  "newNonce": "https://acme-v02.api.letsencrypt.org/acme/new-nonce",
  "newOrder": "https://acme-v02.api.letsencrypt.org/acme/new-order",
  "revokeCert": "https://acme-v02.api.letsencrypt.org/acme/revoke-cert"
}
```

newAccount:  字段是用于上传 account key 时候调用，
newNonce:    业务流程（如申请、撤销） 时候需要预请求获取一次性令牌的一个接口，
newOrder:    用于下单的接口，
revokeCert:  用于撤销证书的接口，
keyChange:   是用于更换 KeyPair 的接口。
4bmRla_FAhE: 特殊的随机的字段，可能随意的出现在 json 的头、中、尾部，作用是为了避免开发者集成客户端时候[1]，采用 定位法取 json 的数据，导致扩展性变差而引入的属性。


2. /acme/new-acct



2.1 HEAD /acme/new-nonce

header: Replay-Nonce 必须提供


2.2 POST /acme/new-acct

```json
{
    "protected": "eyJub25jZSI6ICJ0LWdpTE85RUFMaE9aM2lCM0gtU1JTRXd4ZV9MQTFrN3JaN2FKa0hnbW5ZIiwgInVybCI6ICJodHRwczovL2FjbWUtdjAyLmFwaS5sZXRzZW5jcnlwdC5vcmcvYWNtZS9uZXctYWNjdCIsICJhbGciOiAiUlMyNTYiLCAiandrIjogeyJlIjogIkFRQUIiLCAia3R5IjogIlJTQSIsICJuIjogInV4b2diWk1tM2hYcTVjM3RnSWtxdHJLeWs2OVpXS09sZm1iSG1TRXJkNE5ITmFTeXdmS2V5MmwzLXRBUHl3TF9pSTZJSEVXNFZpWVQ3aGE0MWg2cTZmR3h3dVJjMnY5QUdlZ1YyNFJNd1pYbFp3aWhMQm4zaDUwN1dYUVlJQTVhdVNKUzZOV2tGMEVSc3gyXzA5NU1KQ3JzSDVIdnFiQ3lvY3Jxd2g5N21IeTJ6Qk40SU9Qbi1CRWdiMnhmOF9Bc2xnelFjVlo1VWp6WnhRLXA1bHl6d09wZkRQWFk4LWJzcUZmazktTXUxa3pEdENUYXM2SjM1ZW0yR2hhSFZOODVpX3RBMkxOZmhiejFiTlRvcHQyd1pDRlJ2OWVNaGx4QndfUmhtQ0FrcG1jVm5MLWE2cENuQ1dEU0w4YVd3Mm1HbjVzdVhFNjA3RFdweVhIdm1XSng1USJ9fQ",
    "payload": "eyJjb250YWN0IjogWyJtYWlsdG86IGRvdHNhZmVAemhpaHUuY29tIl0sICJ0ZXJtc09mU2VydmljZUFncmVlZCI6IHRydWV9",
    "signature": "L6aMUAsokgo3BHlxJwvgOAljQZ_bFwUbT2Gc2e3BW1DFXtzmH2qaUg_o-gCyzIqJmc-VAzLG_ipPN2dZRgUl1ARTHHRuniSgItLbSluxoXf138DlW6CJPI-Ma5jYSKoANRnQldAJg_6R07fTfgkzGDp3H2ldZQ-I7QsWmEzp9eRgmpcjIocZJc_AcfbvkK_EPOCn2YSDgDbMUzhwztU3pz5NgQ5tlLO2S_FSQlLlNzv6j4fEU14rR3RQKvl3HmZUAniqXGvsiiiiaYimk815koNei56ZeKTdW0T_lLL1JY3GZbHmK_edzBfKDcaItBdPnnL8FW-bbMB53bMKydpDWw"
}
```


protected base64 decode:

```json
{
    "nonce": "t-giLO9EALhOZ3iB3H-SRSEwxe_LA1k7rZ7aJkHgmnY",
    "url": "https://acme-v02.api.letsencrypt.org/acme/new-acct",
    "alg": "RS256",
    "jwk": {
        "e": "AQAB",
        "kty": "RSA",
        "n": "uxogbZMm3hXq5c3tgIkqtrKyk69ZWKOlfmbHmSErd4NHNaSywfKey2l3-tAPywL_iI6IHEW4ViYT7ha41h6q6fGxwuRc2v9AGegV24RMwZXlZwihLBn3h507WXQYIA5auSJS6NWkF0ERsx2_095MJCrsH5HvqbCyocrqwh97mHy2zBN4IOPn-BEgb2xf8_AslgzQcVZ5UjzZxQ-p5lyzwOpfDPXY8-bsqFfk9-Mu1kzDtCTas6J35em2GhaHVN85i_tA2LNfhbz1bNTopt2wZCFRv9eMhlxBw_RhmCAkpmcVnL-a6pCnCWDSL8aWw2mGn5suXE607DWpyXHvmWJx5Q"
    }
}
```

nonce: /acme/new-nonce所返回的值.
url: 本接口地址
alg: 为签名算法
jwk: 账户公钥数据
kty: key type的意思
e 和 n、kty 可以确定一份公钥



payload base64 decode:

```json
{
    "contact": [
        "mailto: dotsafe@zhihu.com"
    ],
    "termsOfServiceAgreed": true
}
```

protected, payload 和 signature 一起请求给 /acme/new-acct, 服务器存储好公钥与邮箱

http response:

```json
{
  "id": 61907904,
  "key": {
    "kty": "RSA",
    "n": "uxogbZMm3hXq5c3tgIkqtrKyk69ZWKOlfmbHmSErd4NHNaSywfKey2l3-tAPywL_iI6IHEW4ViYT7ha41h6q6fGxwuRc2v9AGegV24RMwZXlZwihLBn3h507WXQYIA5auSJS6NWkF0ERsx2_095MJCrsH5HvqbCyocrqwh97mHy2zBN4IOPn-BEgb2xf8_AslgzQcVZ5UjzZxQ-p5lyzwOpfDPXY8-bsqFfk9-Mu1kzDtCTas6J35em2GhaHVN85i_tA2LNfhbz1bNTopt2wZCFRv9eMhlxBw_RhmCAkpmcVnL-a6pCnCWDSL8aWw2mGn5suXE607DWpyXHvmWJx5Q",
    "e": "AQAB"
  },
  "contact": [
    "mailto: dotsafe@zhihu.com"
  ],
  "initialIp": "00e::3::a11::89::4d",
  "createdAt": "2019-07-24T11:27:45Z",
  "status": "valid"
}
```

创建账户完成


3. /acme/new-order

申请证书的 /acme/new-order 请求之前也需要获取一次性令牌 (如果 new-acct 和 new-order 是一起的,只需要获取一次).



3.1 POST /acme/new-order

```json
{
    "protected": "eyJub25jZSI6ICJYMXAydTR0akVyb29yenNDLVNZZFZrN1hvMFZIaktnSUI5OGVxOGlzLUJzIiwgInVybCI6ICJodHRwczovL2FjbWUtdjAyLmFwaS5sZXRzZW5jcnlwdC5vcmcvYWNtZS9uZXctb3JkZXIiLCAiYWxnIjogIlJTMjU2IiwgImtpZCI6ICJodHRwczovL2FjbWUtdjAyLmFwaS5sZXRzZW5jcnlwdC5vcmcvYWNtZS9hY2N0LzYxOTI0NzQ0In0",
    "payload": "eyJpZGVudGlmaWVycyI6IFt7InR5cGUiOiJkbnMiLCJ2YWx1ZSI6Ind3dy56aGlodS5jb20ifV19",
    "signature": "bzW7CyT2ln1Ms17vfc8kgDgy9yNlwFWFshEUVrFFBUJVWVROvlBjG6CFDXldCGkPvTJCJ1mRvXoMsk3uQ5FhTNaMpBsv7QVL9N03LuLZax1RE_2eyjj5ATMGq8wzXqgAjzC7ueJfcHYuFbuR33NsDXMWZmpEqKi4cf8h-GBDWnuKKUUN0-N-2aQo_JEbRrk5dMbZ1ACJ2Evy46y7-0jCmrHLgnGK7nPK6fmp__6WfNOK_PzjFj3t65NXqio68w2IqXy4CONqx0miqLdgyZfk-HGCN9c4ACNGEhj_mJ9pUEiABaQJagcOqkIADEESg6tIaw9pj0HjqxxlpHTzgrGwxQ"
}
```


base64 decode
```json
{
    "protected": {
        "nonce": "X1p2u4tjEroorzsC-SYdVk7Xo0VHjKgIB98eq8is-Bs",
        "url": "https://acme-v02.api.letsencrypt.org/acme/new-order",
        "alg": "RS256",
        "kid": "https://acme-v02.api.letsencrypt.org/acme/acct/61924744"
    },
    "payload": {
        "identifiers": [
            {
                "type": "dns",
                "value": "www.zhihu.com"
            }
        ]
    },
    "signature": "bzW7CyT2ln1Ms17vfc8kgDgy9yNlwFWFshEUVrFFBUJVWVROvlBjG6CFDXldCGkPvTJCJ1mRvXoMsk3uQ5FhTNaMpBsv7QVL9N03LuLZax1RE_2eyjj5ATMGq8wzXqgAjzC7ueJfcHYuFbuR33NsDXMWZmpEqKi4cf8h-GBDWnuKKUUN0-N-2aQo_JEbRrk5dMbZ1ACJ2Evy46y7-0jCmrHLgnGK7nPK6fmp__6WfNOK_PzjFj3t65NXqio68w2IqXy4CONqx0miqLdgyZfk-HGCN9c4ACNGEhj_mJ9pUEiABaQJagcOqkIADEESg6tIaw9pj0HjqxxlpHTzgrGwxQ"
}
```


protected: 携带的是公钥,
signature: 用私钥将公钥
payload:   取摘要得到的字串
    identifiers: 数组，每一个包含的对象，都有两个属性type和value
        type: 一般为dns
        value: 域名类型为IP时，此type可为ipv4。



http response:

http code: 201
status: pending 表示验证中。可能出现的值分别为 pending/ready/processing/valid/invalid
expires: 订单失效时间。有服务商或CA决定
identifiers同上为数组，每一个包含的对象，都有两个属性type和value
authorizations表示此订单需要依次完成的授权验证资源（Auth-Z）的链接，
authorizations不允许为空数组 （必须至少有一个流程）
finalize授权验证完成后，调用finalize接口签发证书（包括CSR也是在这一步提交的）


```json
{
    "status": "pending",
    "expires": "2019-07-31T16:01:29.887180422Z",
    "identifiers": [
        {
            "type": "dns",
            "value": "www.zhihu.com"
        }
    ],
    "authorizations": [
        "https://acme-v02.api.letsencrypt.org/acme/authz/iJTvo1TYkd16sOJgLZQhCtDx7V3siojiH0tdFVzZPtI"
    ],
    "finalize": "https://acme-v02.api.letsencrypt.org/acme/finalize/61924744/774083503"
}
```


3.2 POST /acme/authz/*

request:

```json
{
    "protected": "eyJub25jZSI6ICJEbk5SR1V3d1NPMDNKRGJ1c25YY3YzTWFISEQ4TUIzY1VNS0x3cU1PSWFVIiwgInVybCI6ICJodHRwczovL2FjbWUtdjAyLmFwaS5sZXRzZW5jcnlwdC5vcmcvYWNtZS9hdXRoei9pSlR2bzFUWWtkMTZzT0pnTFpRaEN0RHg3VjNzaW9qaUgwdGRGVnpaUHRJIiwgImFsZyI6ICJSUzI1NiIsICJraWQiOiAiaHR0cHM6Ly9hY21lLXYwMi5hcGkubGV0c2VuY3J5cHQub3JnL2FjbWUvYWNjdC82MTkyNDc0NCJ9",
    "payload": "",
    "signature": "ecIpPr-WmaRcfui_SmDFD9-coaaZkP0zbZltE94wkHfI8myoVz1ojgrDb0vJB3iBc16-s1f099Ags3T1Ms2Pzu7_eLBY6T7V4zjoB7hDGvQrFqkVdHbUYHzA64Bv_ULjt-PXKLvAucBIZx5Vzkd9pFeEiF9jTnrAOtUyN1x0DkKFOyS-hxPD43C9l30U7xDfwbU9bHtRyKgoS79hP3GdQaVTFHkRNi1K8FHRC9kaqmI4pJ49OwW8pVoNPNnTjyag6TY121Mp1W2G-Sm-UYP1rjNv7S9OhfP0SzBUy5C8RKAj156XdUFWGOwxq1qH2O5JG7KfzpWwdTDfPTIyjo21vQ"
}
```

response:

```json
{
  "identifier": {
    "type": "dns",
    "value": "www.zhihu.com"
  },
  "status": "pending",
  "expires": "2019-07-31T16:01:29Z",
  "challenges": [
    {
      "type": "http-01",
      "status": "pending",
      "url": "https://acme-v02.api.letsencrypt.org/acme/challenge/iJTvo1TYkd16sOJgLZQhCtDx7V3siojiH0tdFVzZPtI/18676978378",
      "token": "NUO_qleY0NGwfTmAX580Tt7oCLpglUh34nmq-6YkU-A"
    },
    {
      "type": "dns-01",
      "status": "pending",
      "url": "https://acme-v02.api.letsencrypt.org/acme/challenge/iJTvo1TYkd16sOJgLZQhCtDx7V3siojiH0tdFVzZPtI/18676978379",
      "token": "m4heiWkRAKPZ2r9_n5kLkS_g6flckRTiVYDqPjQtFxk"
    },
    {
      "type": "tls-alpn-01",
      "status": "pending",
      "url": "https://acme-v02.api.letsencrypt.org/acme/challenge/iJTvo1TYkd16sOJgLZQhCtDx7V3siojiH0tdFVzZPtI/18676978381",
      "token": "-VZddxtXSpiDgwCoWbAF8FfnSVnn9a_b_XwfXxvnUT8"
    }
  ]
}
```


identifier不再赘述，其为identifiers 中的一项
status表示的是该 identifier的验证通过与否状态
expires 和 order 的失效时间原则上应保持一致
challenges为重要字段，其为数组，返回的数目几乎都包含 type =http-01、type=dns-01、type=tls-alpn-01三条（若申请的域名符号为IP，则没有dns-01），
challenge下的type表示验证方式，目前包含http-01、dns-01和tls-alpn-01 。
status表示此种验证方式的状态
url为这条challenge手动呼唤 CA 供应商发起域名控制器校验的请求地址
token为执行此验证规则的凭据（你申请某域名，跟别人申请同一域名的证书所需要填写/上传的凭据是不同的）


域名所有权验证:
    1. 若选择--webroot选项 （不同acme客户端可能参数不一样，这里是以acme.sh举例），其会在 .well-known/pki-validation/路径下创建一个文件，名为token，内容为token
    2. 若选择--dns {dns_provider}选项，其会在 dns_provider调用API，将控制器在你名下的具体域名添加一项 _acme-challenge.你的域名（若通配符域名，则为_acme-challenge.一级域名） ，类型txt，值为token







author
------

qx@startops.com.cn
