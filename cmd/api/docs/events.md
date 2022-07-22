#### Starting a Tournament Match

```http
  POST /tournaments/{id}/matches/{match_id}/events/start
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Goal a Tournament Match

```http
  POST /tournaments/{id}/matches/{match_id}/events/goal
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter    | Type     | Description                 |
| :----------- | :------- | :-------------------------- |
| `team_score` | `string` | **Required**. Team score id |
| `player`     | `string` | **Required**. Player id     |
| `minute`     | `int`    | **Required**. Goal minute   |

#### Setting Halftime for a Tournament Match

```http
  POST /tournaments/{id}/matches/{match_id}/events/halftime
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Adding an extratime for a Tournament Match

```http
  POST /tournaments/{id}/matches/{match_id}/events/extratime
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter   | Type  | Description                     |
| :---------- | :---- | :------------------------------ |
| `extratime` | `int` | **Required**. Extratime minutes |

#### Adding a warning for a Tournament Match

```http
  POST /tournaments/{id}/matches/{match_id}/events/warning
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter | Type      | Description                                        |
| :-------- | :-------- | :------------------------------------------------- |
| `team`    | `string`  | **Required**. Team id                              |
| `player`  | `string`  | **Required**. Player id                            |
| `warning` | `warning` | **Required**. Warning type - [RedCard, YellowCard] |
| `minute`  | `int`     | **Required**. Warning minute                       |

#### Substitution players for a Tournament Match

```http
  POST /tournaments/{id}/matches/{match_id}/events/substitution
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter    | Type     | Description                       |
| :----------- | :------- | :-------------------------------- |
| `team`       | `string` | **Required**. Team id             |
| `player_out` | `string` | **Required**. Player out id       |
| `player_in`  | `string` | **Required**. Player in id        |
| `minute`     | `int`    | **Required**. Substitution minute |

#### Finishing a Tournament Match

```http
  POST /tournaments/{id}/matches/{match_id}/events/finish
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |
