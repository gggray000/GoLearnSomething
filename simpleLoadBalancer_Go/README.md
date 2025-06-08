## Usage

Build<br>
<code>go build -o loadbalancer main.go
</code>

Start a backend server at port 3031:<br>
<code>go run backend.go</code>

Start load balancer at port 3030:<br>
<code>./loadbalancer --backends=http://localhost:3031 --port=3030 </code>

### Acknowledgement
https://github.com/kasvith/simplelb/