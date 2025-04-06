# European Go Federation Rating System

This is repository contains a Rest API with and endpoint to calculate the [EGF Rating System](https://europeangodatabase.eu/EGD/EGF_rating_system.php).

## Endpoints

API Endpoints are expose on port `9000`.

### `new_route`

Calculate new rating of player given both player ratings and the result of the game. This only calculate the `new_rating`and `gor_change` for the `player_rating`!

Input

```json
{
  "player_rating": 2674.564,
  "opponent_rating": 2611.051,
  "result": 0.0
}
```

Response

```json
{
  "new_rating": 2670.45617061712,
  "gor_change": -4.107829382880027,
  "expected_result": 0.6630749860799611,
  "con": 6.197900613388789,
  "bonus": 0.0018434800675773823,
  "beta": -45.06914304568637
}
```

## System description

Ratings are updated by: `r' = r + con * (Sa - Se) + bonus`

`r` is the old EGD rating (GoR) of the player  
`r'` is the new EGD rating of the player  
`Sa` is the actual game result (1.0 = win, 0.5 = jigo, 0.0 = loss)  
`Se` is the expected game result as a winning probability (1.0 = 100%, 0.5 = 50%, 0.0 = 0%). See further below for its computation.  
`con` is a factor that determines rating volatility (similar to K in regular Elo rating systems): `con = ((3300 - r) / 200)^1.6`  
`bonus` A term included to counter rating deflation: `bonus = ln(1 + exp((2300 - r) / 80)) / 5`

`Se` is computed by the Bradley-Terry formula: `Se = 1 / (1 + exp(β(r2) - β(r1)))`  
`r1` is the EGD rating of the player  
`r2` is the EGD rating of the opponent  
`β` is a mapping function for EGD ratings: `β = -7 * ln(3300 - r)`


