import os

os.environ["LOTUSMARKET_QUIET"] = "1"

import json
from lotusmarket.fetchers.vps import parse_vps_response, _parse_pipe_vol
from lotusmarket.fetchers.entrade import parse_entrade_response


def test_parse_vps_valid():
    raw = json.dumps(
        [
            {
                "sym": "ACB",
                "lastPrice": 26.9,
                "r": 26.5,
                "openPrice": 26.6,
                "highPrice": 27.0,
                "lowPrice": 26.3,
                "lot": 1234567,
                "fBVol": 50000,
                "fSVolume": 30000,
                "c": 28.3,
                "f": 24.7,
                "g1": "26.8|100570|d",
                "g2": "27.0|80000|d",
            }
        ]
    )
    stocks = parse_vps_response(raw.encode())
    assert len(stocks) == 1
    s = stocks[0]
    assert s.ticker == "ACB"
    assert s.close == 26900
    assert s.foreign_net_vol == 20000


def test_parse_vps_empty():
    stocks = parse_vps_response(b"[]")
    assert len(stocks) == 0


def test_parse_entrade_valid():
    raw = json.dumps(
        {
            "t": [1709510400, 1709596800],
            "o": [26.5, 27.0],
            "h": [27.0, 27.5],
            "l": [26.0, 26.5],
            "c": [26.8, 27.2],
            "v": [1000000, 1500000],
        }
    )
    stocks = parse_entrade_response(raw.encode(), "ACB")
    assert len(stocks) == 2
    assert stocks[0].close == 26800


def test_parse_pipe_vol():
    assert _parse_pipe_vol("22.5|100570|d") == 100570
    assert _parse_pipe_vol("") == 0
