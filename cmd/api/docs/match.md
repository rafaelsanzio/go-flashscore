#### Creating a Tournament Match

```http
  POST /tournaments/{id}/matches
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

| Parameter       | Type          | Description                |
| :-------------- | :------------ | :------------------------- |
| `home_team`     | `string`      | **Required**. Home team id |
| `away_team`     | `string`      | **Required**. Away team id |
| `date_of_match` | `string date` | **Required**. Match date   |
| `time_of_match` | `string time` | **Required**. Match time   |

#### Deleting a Tournament Match

```http
  DELETE /tournaments/{id}/matches/{match_id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Getting a Tournament Match

```http
  GET /tournaments/{id}/matches/{match_id}
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |

#### Listing all Tournaments Matches

```http
  GET /tournaments{id}/matches
```

| Header  | Type     | Description                |
| :------ | :------- | :------------------------- |
| `Token` | `Bearer` | **Required**. Your API key |
