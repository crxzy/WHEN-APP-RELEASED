import requests

url = 'https://app.market.xiaomi.com/apm/search?channel=market_100_1_android&clientId=8459b0e8d3067ea452730229756b439f&co=CN&densityScaleFactor=1.6875&imei=06222f21cea9dd3e56569d5b493f11a1&keyword=网易亲时光&la=zh&marketVersion=147&model=MuMu&os=eng.root.20200623.095831&page=0&ref=input&resolution=810*1440&sdk=23'
data = requests.get(url).json()['listApp']
# print(resp)
# print(resp['listApp'])



for app in data:
        if app["displayName"] == "网易亲时光":
           a = app["versionName"]
           print(a)
           b = app["updateTime"]
           print(b)
           updateTime=1609992964524