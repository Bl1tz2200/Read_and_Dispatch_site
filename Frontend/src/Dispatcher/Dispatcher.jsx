import dispatchLogo from '../assets/dispatch_book.svg'
import dispatchAnimation from '../assets/dispatch.gif'
import './Dispatcher.css'
import { Fragment, useEffect, useState } from 'react'
import { useAuth } from '../Providers/AuthProvider'
import { API_ADRESS } from "../main.jsx"
import { ShowMiniInfo } from '../PublicFunctions/ShowMiniInfo.jsx'

export function Dispatcher() { // Main page where users can see all dispatches
  const { token, setToken } = useAuth();
  const [dispatchIds, setdispatchIds] = useState([])
  const [userName, setUserName] = useState(null)

  async function getUsername(token) { // Request to backend to get username
    const result = await fetch(`${API_ADRESS}/auth`, {
      headers:
      {
        "Authorization": token,
      }
    }).then(response => {

      if (response.status == 202) {

        return response.json()

      } else if (response.status == 401) {

        return response.json()

      } else {

        return null

      }

    }).then(responseMessage => {
      if (responseMessage) { // If user is authed

        if (responseMessage.UserName) { // If no errors

          return responseMessage.UserName

        } else { // If user auth is invalid

          setToken()
          window.location.reload()

        }

      }

    })

    return await result
  }

  async function getDispatches() { // Request to backend to get all dispatches it
    const ids = await fetch(`${API_ADRESS}/getDispatches`).then(response => {

      if (response.status == 202) {

        return response.json()

      } else if (response.status == 401) {

        alert("Something went wrong while authorization! Try to log in again")
        setToken()
        window.location.pathname = '/login'

      } else {

        return null

      }

    }).then(ids => {
      if (ids) { // If no errors

        setdispatchIds(ids)

      } else { // If error appeared

        alert("Something went wrong on a server! Try again later")
        window.location.reload()

      }
    })

  }

  useEffect(() => { // Functions that will be run when page was loaded
    getDispatches()

    if (token) { // If user is authed
      getUsername(token).then(
        username => {
          setUserName(username)
        })
    }
  }, [])

  return (
    <>
      <div className="header">
        <a href="/dispatch"><img src={dispatchLogo} className="logo" alt="Dispatch"
          onMouseOver={mySrc => { mySrc.currentTarget.src = dispatchAnimation; mySrc.currentTarget.className = "animatedLogo" }}
          onMouseOut={mySrc => { mySrc.currentTarget.src = dispatchLogo; mySrc.currentTarget.className = "logo" }} /></a>
        <input type='text' placeholder='Search in all works...' id="searchAll" onKeyDown={(enteredKey) => {
          if (enteredKey.key == "Enter"){

            var enteredText = document.getElementById("searchAll").value

            dispatchIds.map(id =>
            {
  
              if(!(`${document.getElementById(`${id}`).textContent}`.includes(enteredText))){
                document.getElementById(`${id}/div`).style.display = "none"
              } else {
                document.getElementById(`${id}/div`).style.display = "inline-flex"
              }

            }
            )
          }
        }} />
        <div className="userName">
          <h1><a href='/user'>{userName ? userName : <>Log into <br />account</>}</a></h1>
        </div>
      </div>
      <div className="articles">
        <ul>
          {
            dispatchIds.map(id =>
              <Fragment key={id}>
                <a href={`/dispatch/${id}`}>< ShowMiniInfo id={id} /></a>
              </Fragment>
            )
          }
        </ul>
      </div>
    </>
  )
}
