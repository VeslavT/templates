import hmac
import base64
import hashlib
from datetime import datetime
from typing import Dict
from functools import wraps
import requests
from requests.exceptions import HTTPError
from urllib.parse import urlparse, urljoin

from flask import current_app, request, make_response, jsonify

API_SECRET = ""


def sign_request(path: str, body: str = "") -> Dict[str, str]:
    """
    :param path: Url for signing
    :param body: Post body data
    :return: dict
         headers
         {"X-Application-Sign": ..., "X-Application-Sign-Date": ...}
    """
    url = urlparse(path)
    if not url.query:
        # sign algorithm for chat request "?"
        path += "?"
    dt = datetime.utcnow().strftime("%Y-%m-%d %H:%M:%S")
    body = path.encode('utf-8') + body.encode('utf-8') + dt.encode('utf-8')
    mac = hmac.new(
        base64.b64decode(API_SECRET),
        body,
        hashlib.sha256).hexdigest()
    return {"X-Application-Sign": mac, "X-Application-Sign-Date": dt}


def check_sign_request(req, secret):
    """
    Check correct signature of request
    Can be used with decorator @sign_required
    :param req: flask Request object.
    :param secret: secret for sign request.
    """
    sign_data = req.headers.get("X-Application-Sign", "")
    sign_date = req.headers.get("X-Application-Sign-Date", "")
    if sign_data == "" or sign_date == "":
        return False

    try:
        dt = datetime.strptime(sign_date, "%Y-%m-%d %H:%M:%S")
    except ValueError:
        current_app.logger.info("Incorrect data format '%s', should be 'YYYY-MM-DD HH:MM:SS'", sign_date)
        return False

    tdelta = datetime.utcnow() - dt
    delta = int(tdelta.total_seconds() / 60)
    if abs(delta) > 15:
        current_app.logger.error("Incorrect(expired) signature")
        return False

    url = urlparse(req.path)
    path = req.path if url.query else req.path + "?"
    body = request.get_data()
    expected_sign = hmac.new(
        base64.b64decode(secret),
        path.encode('utf-8') + body + sign_date.encode('utf-8'),
        hashlib.sha256).hexdigest()

    if expected_sign != sign_data:
        current_app.logger.error("Invalid signature")
        return False
    return True


def sign_required(f):
    @wraps(f)
    def decorated_view(*args, **kwargs):
        if not current_app.debug and not check_sign_request(request, API_SECRET):
            return make_response(jsonify("Correct signature required"), 403)
        return f(*args, **kwargs)
    return decorated_view


def get_users_request(domain, path):
    """
       Example of usage sign of request
       :param domain: requested domain.
       :param path: requested URI.
       """
    headers = sign_request(path)
    headers["Content-Type"] = "application/json"
    r = requests.get(urljoin(domain, path), headers=headers)
    if r.status_code != 200:
        raise HTTPError(u'%s Error: %s for url: %s' % (r.status_code, r.reason, r.url), response=r)
    return r.json()
