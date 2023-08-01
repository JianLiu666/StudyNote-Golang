package main

var ecSiteIdPool = []string{
	"1",
	"21729934",
	"21729963",
	"21729991",
	"21730024",
	"21730063",
	"21730086",
	"21730112",
	"21730139",
	"21730184",
	"21730208",
	"21730232",
	"3319009",
	"5044073",
	"63638",
}

var baseWagerRecord = `
{
  "startTimestamp": 1677934262839,
  "endTimestamp": 1677934506671,
  "playerId": "1777614979",
  "nickname": "lucas123456",
  "ecName": "testEcSite",
  "rawEcName": "auth-1",
  "recordType": "Game",
  "gameName": "深海历险",
  "gameRecordId": "GR-1677934506671-23-56172",
  "sessionRecordId": "WGR-1677934506671-23-56172-1777614979",
  "seq": 0,
  "kkBeforBalance": 44161857200,
  "kkAfterBalance": 44607357200,
  "beforeBetBalance": 44161857200,
  "kkDepositAmount": 0,
  "kkDrawAmount": 0,
  "depositAmount": 0,
  "drawAmount": 0,
  "betAmount": 15000000,
  "betBalance": 44607357200,
  "validBetAmount": 15000000,
  "memberJuiceAmount": 0,
  "memberHostJuiceAmount": 0,
  "profit": 445500000,
  "payoutAmount": 460500000,
  "balance": 44607357200,
  "ecUserId": "lucas123456",
  "ecSiteId": "1",
  "gameId": "4",
  "gameType": 1,
  "scoreType": 1
}`

var baseGameRecord = `
{
  "ecSiteId": "1",
  "gameType": 1,
  "mergedGameRecordId": "",
  "gameId": "4",
  "themeId": "1",
  "roomId": "624",
  "memberCount": 1,
  "memberData": [
    {
      "ecUserId": "lucas123456",
      "playerId": "1777614979",
      "sessionId": "SR-1677933048541-0-9",
      "wagerId": "WGR-1677934506671-23-56172-1777614979",
      "nickname": "lucas123456",
      "icon": 0,
      "iconFrame": 0,
      "kkBeforBalance": 44161857200,
      "kkAfterBalance": 44607357200,
      "beforeBetBalance": 44161857200,
      "betBalance": 44607357200,
      "betAmount": 15000000,
      "validBetAmount": 15000000,
      "memberIncome": 460500000,
      "memberOutcome": 15000000,
      "memberJuiceAmount": 0,
      "memberHostJuiceAmount": 0,
      "profit": -445500000,
      "payoutAmount": 460500000,
      "score": 0,
      "tagType": 2
    }
  ],
  "memberIncomes": 460500000,
  "memberOutcomes": 15000000,
  "memberJuiceAmounts": 0,
  "memberHostJuiceAmounts": 0,
  "botCount": 0,
  "botIncome": 0,
  "botOutcome": 0,
  "botJuiceAmount": 0,
  "juiceRatio": 0,
  "history": [
    {
      "AllReels": [
        [
          13,
          20,
          5
        ],
        [
          11,
          20,
          1
        ],
        [
          5,
          12,
          3
        ],
        [
          2,
          20,
          12
        ],
        [
          3,
          14,
          5
        ]
      ],
      "RandomNumList": [
        27,
        29,
        0,
        32,
        0
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 75000000,
      "EliminateList": [
        {
          "WinPos": [
            [
              0,
              1,
              0
            ],
            [
              0,
              1,
              0
            ],
            [
              0,
              0,
              0
            ],
            [
              0,
              1,
              0
            ],
            [
              0,
              0,
              0
            ]
          ],
          "IsEliminate": false,
          "NextReels": [
            [
              13,
              20,
              5
            ],
            [
              11,
              20,
              1
            ],
            [
              5,
              12,
              3
            ],
            [
              2,
              20,
              12
            ],
            [
              3,
              14,
              5
            ]
          ],
          "WinBonus": 75000000,
          "WinLineList": [
            {
              "SpecialMultiple": 1,
              "X": 1,
              "WinNum": 20,
              "WinBonus": 75000000,
              "LineNum": 3
            }
          ]
        }
      ]
    },
    {
      "AllReels": [
        [
          11,
          14,
          13
        ],
        [
          3,
          1,
          12
        ],
        [
          4,
          20,
          14
        ],
        [
          12,
          2,
          13
        ],
        [
          5,
          11,
          1
        ]
      ],
      "RandomNumList": [
        56,
        16,
        33,
        17,
        49
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          12,
          11,
          4
        ],
        [
          12,
          13,
          1
        ],
        [
          14,
          11,
          5
        ],
        [
          13,
          2,
          1
        ],
        [
          11,
          4,
          13
        ]
      ],
      "RandomNumList": [
        5,
        11,
        81,
        43,
        57
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          5,
          13,
          4
        ],
        [
          5,
          12,
          14
        ],
        [
          2,
          13,
          4
        ],
        [
          13,
          1,
          14
        ],
        [
          14,
          1,
          50
        ]
      ],
      "RandomNumList": [
        51,
        54,
        43,
        19,
        60
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          5,
          3,
          12
        ],
        [
          11,
          1,
          12
        ],
        [
          2,
          1,
          3
        ],
        [
          5,
          13,
          11
        ],
        [
          2,
          14,
          14
        ]
      ],
      "RandomNumList": [
        28,
        35,
        39,
        13,
        78
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          14,
          11,
          2
        ],
        [
          1,
          13,
          12
        ],
        [
          14,
          11,
          4
        ],
        [
          2,
          13,
          1
        ],
        [
          20,
          11,
          4
        ]
      ],
      "RandomNumList": [
        8,
        31,
        77,
        18,
        56
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          14,
          4,
          11
        ],
        [
          13,
          1,
          12
        ],
        [
          12,
          14,
          11
        ],
        [
          5,
          12,
          11
        ],
        [
          4,
          11,
          5
        ]
      ],
      "RandomNumList": [
        31,
        12,
        26,
        35,
        30
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          13,
          20,
          2
        ],
        [
          4,
          1,
          13
        ],
        [
          50,
          50,
          50
        ],
        [
          14,
          1,
          12
        ],
        [
          13,
          2,
          4
        ]
      ],
      "RandomNumList": [
        58,
        69,
        9,
        2,
        28
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 40500000,
      "EliminateList": [
        {
          "WinPos": [
            [
              1,
              0,
              0
            ],
            [
              0,
              0,
              1
            ],
            [
              1,
              1,
              1
            ],
            [
              0,
              0,
              0
            ],
            [
              0,
              0,
              0
            ]
          ],
          "IsEliminate": true,
          "NextReels": [
            [
              14,
              20,
              2
            ],
            [
              14,
              4,
              1
            ],
            [
              50,
              50,
              50
            ],
            [
              14,
              1,
              12
            ],
            [
              13,
              2,
              4
            ]
          ],
          "WinBonus": 4500000,
          "WinLineList": [
            {
              "SpecialMultiple": 1,
              "X": 3,
              "WinNum": 13,
              "WinBonus": 4500000,
              "LineNum": 3
            }
          ]
        },
        {
          "WinPos": [
            [
              1,
              0,
              0
            ],
            [
              1,
              0,
              0
            ],
            [
              1,
              1,
              1
            ],
            [
              1,
              0,
              0
            ],
            [
              0,
              0,
              0
            ]
          ],
          "IsEliminate": true,
          "NextReels": [
            [
              11,
              20,
              2
            ],
            [
              11,
              4,
              1
            ],
            [
              14,
              1,
              50
            ],
            [
              12,
              1,
              12
            ],
            [
              13,
              2,
              4
            ]
          ],
          "WinBonus": 27000000,
          "WinLineList": [
            {
              "SpecialMultiple": 2,
              "X": 3,
              "WinNum": 14,
              "WinBonus": 27000000,
              "LineNum": 4
            }
          ]
        },
        {
          "WinPos": [
            [
              1,
              0,
              0
            ],
            [
              1,
              0,
              0
            ],
            [
              0,
              0,
              1
            ],
            [
              0,
              0,
              0
            ],
            [
              0,
              0,
              0
            ]
          ],
          "IsEliminate": true,
          "NextReels": [
            [
              14,
              20,
              2
            ],
            [
              13,
              4,
              1
            ],
            [
              3,
              14,
              1
            ],
            [
              12,
              1,
              12
            ],
            [
              13,
              2,
              4
            ]
          ],
          "WinBonus": 9000000,
          "WinLineList": [
            {
              "SpecialMultiple": 3,
              "X": 1,
              "WinNum": 11,
              "WinBonus": 9000000,
              "LineNum": 3
            }
          ]
        }
      ]
    },
    {
      "AllReels": [
        [
          11,
          5,
          4
        ],
        [
          20,
          1,
          13
        ],
        [
          14,
          11,
          4
        ],
        [
          1,
          12,
          3
        ],
        [
          14,
          4,
          13
        ]
      ],
      "RandomNumList": [
        33,
        30,
        27,
        66,
        43
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          4,
          11,
          14
        ],
        [
          13,
          12,
          1
        ],
        [
          14,
          13,
          3
        ],
        [
          12,
          14,
          1
        ],
        [
          1,
          2,
          12
        ]
      ],
      "RandomNumList": [
        35,
        32,
        74,
        1,
        51
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          5,
          4,
          11
        ],
        [
          12,
          13,
          5
        ],
        [
          3,
          13,
          4
        ],
        [
          5,
          12,
          11
        ],
        [
          11,
          2,
          14
        ]
      ],
      "RandomNumList": [
        34,
        8,
        63,
        35,
        37
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          20,
          2,
          5
        ],
        [
          5,
          12,
          13
        ],
        [
          11,
          14,
          5
        ],
        [
          50,
          50,
          13
        ],
        [
          2,
          13,
          5
        ]
      ],
      "RandomNumList": [
        26,
        10,
        87,
        8,
        11
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 90000000,
      "EliminateList": [
        {
          "WinPos": [
            [
              0,
              0,
              1
            ],
            [
              1,
              0,
              0
            ],
            [
              0,
              0,
              1
            ],
            [
              1,
              1,
              0
            ],
            [
              0,
              0,
              1
            ]
          ],
          "IsEliminate": true,
          "NextReels": [
            [
              13,
              20,
              2
            ],
            [
              13,
              12,
              13
            ],
            [
              2,
              11,
              14
            ],
            [
              50,
              50,
              13
            ],
            [
              5,
              2,
              13
            ]
          ],
          "WinBonus": 90000000,
          "WinLineList": [
            {
              "SpecialMultiple": 1,
              "X": 2,
              "WinNum": 5,
              "WinBonus": 90000000,
              "LineNum": 5
            }
          ]
        }
      ]
    },
    {
      "AllReels": [
        [
          12,
          14,
          4
        ],
        [
          12,
          3,
          11
        ],
        [
          50,
          14,
          5
        ],
        [
          50,
          13,
          12
        ],
        [
          3,
          14,
          20
        ]
      ],
      "RandomNumList": [
        30,
        39,
        16,
        52,
        34
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 255000000,
      "EliminateList": [
        {
          "WinPos": [
            [
              1,
              0,
              0
            ],
            [
              1,
              0,
              0
            ],
            [
              1,
              0,
              0
            ],
            [
              1,
              0,
              1
            ],
            [
              0,
              0,
              0
            ]
          ],
          "IsEliminate": true,
          "NextReels": [
            [
              3,
              14,
              4
            ],
            [
              2,
              3,
              11
            ],
            [
              50,
              14,
              5
            ],
            [
              50,
              50,
              13
            ],
            [
              3,
              14,
              20
            ]
          ],
          "WinBonus": 15000000,
          "WinLineList": [
            {
              "SpecialMultiple": 1,
              "X": 2,
              "WinNum": 12,
              "WinBonus": 15000000,
              "LineNum": 4
            }
          ]
        },
        {
          "WinPos": [
            [
              1,
              0,
              0
            ],
            [
              0,
              1,
              0
            ],
            [
              1,
              0,
              0
            ],
            [
              1,
              1,
              0
            ],
            [
              1,
              0,
              0
            ]
          ],
          "IsEliminate": true,
          "NextReels": [
            [
              5,
              14,
              4
            ],
            [
              12,
              2,
              11
            ],
            [
              50,
              14,
              5
            ],
            [
              50,
              50,
              13
            ],
            [
              12,
              14,
              20
            ]
          ],
          "WinBonus": 240000000,
          "WinLineList": [
            {
              "SpecialMultiple": 2,
              "X": 2,
              "WinNum": 3,
              "WinBonus": 240000000,
              "LineNum": 5
            }
          ]
        }
      ]
    },
    {
      "AllReels": [
        [
          14,
          5,
          12
        ],
        [
          13,
          5,
          12
        ],
        [
          1,
          4,
          20
        ],
        [
          1,
          12,
          50
        ],
        [
          12,
          2,
          13
        ]
      ],
      "RandomNumList": [
        3,
        9,
        32,
        3,
        14
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          11,
          2,
          3
        ],
        [
          1,
          13,
          12
        ],
        [
          12,
          3,
          14
        ],
        [
          12,
          3,
          11
        ],
        [
          2,
          14,
          14
        ]
      ],
      "RandomNumList": [
        9,
        31,
        56,
        54,
        78
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    },
    {
      "AllReels": [
        [
          13,
          20,
          2
        ],
        [
          4,
          14,
          3
        ],
        [
          4,
          14,
          2
        ],
        [
          2,
          12,
          5
        ],
        [
          4,
          5,
          14
        ]
      ],
      "RandomNumList": [
        58,
        42,
        65,
        22,
        71
      ],
      "LongWildIndex": 5,
      "TotalWinBonus": 0,
      "EliminateList": []
    }
  ],
  "startTimestamp": 1677934262839,
  "endTimestamp": 1677934506671,
  "scoreType": 1,
  "signature": "MEQCIDNQcF8ZQk2dH76mEGXShxnfYEO+NHEXgYd1C27H3/P2AiAJPLdidd5wC1mrZ371KMgDx1DjdYQMY8vQijxB1xz/4w=="
}`
