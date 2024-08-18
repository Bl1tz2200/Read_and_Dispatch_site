<div>
  <h1>Read_and_Dispatch_site</h1>
  <p>I made site where you can publist your photos with description. For an API I used golang Gin. I made account registration, log in and authentication too</p>
</div>
<br>
<div>
  <h1>How to start using?</h1>
  <p>
    Be sure, that you have already installed <i>MySql DB</i> server. <br>
    It should work on the default port <strong>3306</strong> <br>
    (you can change port, but then you should change port inside <b>main.go</b>) <br>
    <br>
    After connecting to the DB enter <b>Backend/DB</b> directory and run <b>Users_DB.sql</b> script to create Users_DB.<br>
    This DB will store User's credentials and ids of dispatches in tables <strong>Users</strong> and <strong>Users_Dispatches</strong><br>
    (All Users have unique UserName and UserEmail. All dispatches have unique ID)<br>
    <br>
    For running backend you should install <i>Golang</i> <br>
    For running frontend you should use <i>Vite.js</i> with <i>React</i> (without <i>Typescript</i>)<br>
    <br>
    To run backend you should enter <b>Backend/Golang</b> directory with <b>main.go</b> and <b>go.mod</b> and then write in console: <br>
    <pre>
       $  go get . 
       $  go run . </pre>
    It will install all <i>Go</i> dependencies and run backend on <strong>http://localhost:8080</strong> <br>
    (you can change it inside <b>main.go</b>, but then you should change it inside frontend's <b>main.jsx</b>) <br>
    <br>
    To run frontend you should enter <b>Frontend</b> directory with <b>package.json</b> and other files, then write in console:
    <pre>
       $  npm install  
       $  npm run dev </pre>
    It will install all <i>React</i> dependencies and run frontend on <strong>http://localhost:5173</strong> <br>
    (you can change server's ip and port inside <b>vite.config.js</b>, but then you should change it inside backend's <b>main.go</b>) <br>
    <br>
    <strong><em>After all you will get working website where you can post and read posts!</em></strong>
    <br>
  </p>
  <br>
  <h3>Notice:</h3>
    <h4>
      <ul>
        <li>
          Password and other User's credentials are sended to backend and stored inside DB without any hashing or encryption. <br>
          It could cause troubles with safety. 
        </li>
        <br>
        <li>
          Function getDispatches() in <b>Frontend/Dispatcher/Dispatcher.jsx</b> get all ids of dispatches from DB. <br>
          It could cause troubles with optimization and page rendering if there would be a lot of dispatches. 
        </li>
      </ul>
    </h4>
</div>
