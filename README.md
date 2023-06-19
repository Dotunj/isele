##  Isele

Isele is a CLI tool I use for sending incoming and outgoing events to my [Convoy](https://github.com/frain-dev/convoy) instance. It uses [Vegeta](https://github.com/tsenart/vegeta) under the hood to make the requests. 


### Usage Manual
```
Usage: isele <command> [command flags]

serve command:
  - host string
        The base URL of your convoy instance
  - apiKey string
        API Key for your project
  - projectID string
        The project's ID
  - endpointID string
        This is only required for an outgoing project
  - maskID string
        This is only required for an incoming project
  - duration int
        Duration in seconds
  - rate int
        Number of events to send per second
  - body string
        Path to request body file. This is optional as there is a default payload to be used
  
  examples:
  // Outgoing Project
  isele serve --apiKey=CO.0pnlMcp1O56hhfK1.aAk6AHThhRMjkEuOA0vqOHjJJdqjW4K4geV75N8KUTK48vtMeuvsuAfyIRjMqs3C --endpointID=01H28MXEGATHC7TH1J3Y28WJSV --rate=3 --projectID=01H28MWP1TQV5AC5QR54BSX8HJ --duration=10 --body=events.json
  
  // Incoming Project
  isele serve --apiKey=CO.kbVgUwh8hBdhQuWJ.iR8zCg4FXCh9B9ws4fmcTKyXaNYqZymzall1HIucKCU4bwwmBBg89jKTXjukz5p6 --rate=3 --projectID=01H28MWP1TQV5AC5QR54BSX8HJ --maskID=s87LH7p2lHXMij1L       
```
## License
The MIT License (MIT). Please see [License File](LICENSE.md) for more information.