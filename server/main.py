import flask
from flask import request
from flask.json import jsonify
import requests
import random

app = flask.Flask(__name__)

common = {
    "Version": "",
    "ReleaseTime": "",
    "Msg": ""
}

@app.route("/ios", methods=['POST'])
def ios():
    url = 'https://itunes.apple.com/lookup?bundleId={}&country=CN&rnd={}'
    ret = common.copy()
    req = request.get_json()
    print(req)
    resp = requests.get(url.format(req['BundleID'], str(random.randint(1, 999999999))))

    if resp.status_code != 200:
        ret['Msg'] = "error"
        return jsonify(ret)
    
    data = resp.json()
    ret['Version'] = data['results'][0]['version']
    ret['ReleaseTime'] = data['results'][0]['currentVersionReleaseDate']
    ret['Msg'] = "ok"
    return jsonify(ret)

@app.route("/huawei", methods=['post'])
def huawei():
    url = 'https://web-drcn.hispace.dbankcloud.cn/uowap/index?method=internal.completeSearchWord&serviceType=20&keyword={}&zone=&locale=zh'
    ret = common.copy()
    req = request.get_json()
    print(req)
    resp = requests.get(url.format(req['Name']))

    if resp.status_code != 200:
        ret['Msg'] = "error"
        return jsonify(ret)

    data = resp.json()
    for app in data["appList"]:
        if app["package"] == req["PackageName"]:
            ret['Version'] = app["version"]
            ret['ReleaseTime'] = app['releaseDate']
            ret['Msg'] = 'ok'
            print(ret)
            return jsonify(ret)

    ret['Msg'] = "not found"
    print(ret)
    return jsonify(ret)


@app.route("/myapp", methods=['post'])
def myapp():
    url = 'https://sj.qq.com/myapp/searchAjax.htm?kw={}&pns=&sid='
    ret = common.copy()
    req = request.get_json()
    print(req)
    resp = requests.get(url.format(req['Name']))

    if resp.status_code != 200:
        ret['Msg'] = "error"
        return jsonify(ret)

    data = resp.json()['obj']['items']
    for app in data:
        if app["pkgName"] == req["PackageName"]:
            ret['Version'] = app["appDetail"]["versionName"]
            ret['ReleaseTime'] = str(app["appDetail"]['apkPublishTime'])
            ret['Msg'] = 'ok'
            print(ret)
            return jsonify(ret)

    ret['Msg'] = "not found"
    print(ret)
    return jsonify(ret)

@app.route("/xiaomi_app", methods=['post'])
def xiaomi_app():
    url = 'https://app.market.xiaomi.com/apm/search?channel=market_100_1_android&clientId=8459b0e8d3067ea452730229756b439f&co=CN&densityScaleFactor=1.6875&imei=06222f21cea9dd3e56569d5b493f11a1&keyword={}&la=zh&marketVersion=147&model=MuMu&os=eng.root.20200623.095831&page=0&ref=input&resolution=810*1440&sdk=23'
    ret = common.copy()
    req = request.get_json()
    print(req)
    resp = requests.get(url.format(req['Name']))

    if resp.status_code != 200:
        ret['Msg'] = "error"
        return jsonify(ret)

    data = resp.json()['listApp']
    for app in data:
        if app["packageName"] == req["PackageName"]:
            ret['Version'] = app["versionName"]
            ret['ReleaseTime'] = str(app["updateTime"])
            ret['Msg'] = 'ok'
            print(ret)
            return jsonify(ret)

    ret['Msg'] = "not found"
    print(ret)
    return jsonify(ret)

@app.route("/baidu", methods=['post'])
def baidu():
    url = 'https://mobile.baidu.com/api/appsearch?word={}'
    ret = common.copy()
    req = request.get_json()
    print(req)
    resp = requests.get(url.format(req['Name']))

    if resp.status_code != 200:
        ret['Msg'] = "error"
        return jsonify(ret)

    data = resp.json()['data']['data']
    for app in data:
        if app["package"] == req["PackageName"]:
            ret['Version'] = app["versionname"]
            ret['ReleaseTime'] = str(app["updatetime"])
            ret['Msg'] = 'ok'
            print(ret)
            return jsonify(ret)

    ret['Msg'] = "not found"
    print(ret)
    return jsonify(ret)   
if __name__ == "__main__":
    app.run(host="0.0.0.0", debug=True)