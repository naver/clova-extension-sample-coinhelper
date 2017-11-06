### 개요 
'코인헬퍼'라는 Clova extension의 REST API 서버 소스 코드입니다. '코인헬퍼' 익스텐션이 어떻게 작동하는지 보시려면, Clova 앱이나 Clova 스피커(WAVE, Friens)에서 '코인헬퍼 시작해줘'라고 말해보시길 바랍니다. 해당 익스텐션을 실행하면 코인헬퍼가 비트코인을 비롯한 가상화폐의 시세를 알려줍니다. 

### 사용환경
'코인헬퍼' Clova extension의 REST API 서버는 Go로 구현되어 있습니다.  Windows, MacOS, Linux 등 golang이 구동 가능한 OS면 실행 가능하며, 구체적인 목록들은 여기를 참고하셔서, 아래 가이드에 따라 Go를 먼저 설치하시길 바랍니다. https://github.com/golang-kr/golang-doc/wiki/%EC%84%A4%EC%B9%98-%EC%8B%9C%EC%9E%91%ED%95%98%EA%B8%B0 

### 설치방법
'코인헬퍼' REST API 서버 소스 코드 설치는 다음과 같이 해주시길 바랍니다.
1) Go 배포판 설치: https://github.com/golang-kr/golang-doc/wiki/%EC%84%A4%EC%B9%98-%EC%8B%9C%EC%9E%91%ED%95%98%EA%B8%B0 
2) 소스코드 다운로드:  # git clone https://github.com/naver/clova-extension-sample-coinhelper.git
3) 소스코드 빌드: # make  
소스코드가 정상 빌드되면 bin 디렉토리 밑에 coinHelper라는 실행파일이 생성됩니다. 

### 사용법 
'코인헬퍼' Clova extension의 REST API 서버는 Clova platform으로부터의 익스텐션 요청에 따라 빗썸, 코인원, 코빗 3곳의 가상화폐(비트코인, 이더리움, 이더리움 클래식, 리플)에 대한 시세를 조회한 결과를 응답을 하도록 되어 있습니다. API 서버를 실행을 하더라도, Clova platform이 보내는 것과 동일한 API 요청을 해주셔야 정확하게 작동하는 점 참고 바랍니다. 실제 서비스를 위해서는 https 기반으로 외부에서 접근 가능한 도메인으로 해주셔야 합니다.
- API 서버 실행: bin/coinHelper 
- API 서버 테스팅: [Postman](https://www.getpostman.com/apps)에서 아래와 같이 json Request를 전송하고 json이 리턴되는지 테스트 해봅니다.
	- URL: http://localhost:10680/currency 
	- 요청 방법: POST 
	- Body: raw ( JSON 선택 ) 
- 요청 예시)
```
{
    "version": "0.1.0",
    "session": {
        "sessionId": "e18bc9c3-0881-4781-b271-b34087cf303b",
        "user": {
            "userId": "dO3pmiTPSZ2YwtHvF7bqeA",
            "accessToken": "5e265830-400e-44c8-8c27-31ab2b4781f8"
        },
        "new": true
    },
    "context": {
        "System": {
            "user": {
                "userId": "dO3pmiTPSZ2YwtHvF7bqeA",
                "accessToken": "5e265830-400e-44c8-8c27-31ab2b4781f8"
            },
            "device": {
                "deviceId": "75ef21a0-9c2b-4f13-bb98-9d3ad7aa74ac"
            }
        }
    },
    "request": {
        "type": "IntentRequest",
        "intent": {
            "name": "AskCoinPriceIntent",
            "slots": {
                "currency": {
                    "name": "currency",
                    "value": "비트코인"
                }
            }
        }
    }
}
```

![image](http://static.naver.net/clova/service/native_extensions/example/coinhelper.png)


### 라이선스
Naver & Line corp.

[LICENSE](https://github.com/naver/clova-extension-sample-coinhelper/blob/github-public/LICENSE)

```
Copyright 2018 NAVER Corp. & LINE Corporation

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
```

