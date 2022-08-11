sessionIds = { "357e5c1f-1143-45c7-9b66-a3cc051f25b6" }

gIndex = 1

function request()
    local headers = { }
    method = "GET"
    headers['Content-Type'] = "application/json"
    headers["Cookie"] = "sessionId=" .. sessionIds[1]
    body = ""

    return wrk.format(method, nil, headers, body)
end

function response(status, headers, body)
    if status ~= 200 then
        print("http status error: ", body)
        wrk.thread:stop()
    elseif (not string.find(body, '"code":0')) then
        print("resp code error: ", body)
        wrk.thread:stop()
    end
end  
