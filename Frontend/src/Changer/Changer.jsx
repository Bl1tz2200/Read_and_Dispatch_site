import { useSearchParams } from 'react-router-dom';
import { API_ADRESS } from '../main'
import './Changer.css'

export function Changer() {
  const [searchParams] = useSearchParams();

  async function resetPassword() { // Request to backend to change password 
    const token = searchParams.get("token") // Get mail JWT token from Query string
    console.log(token)

    var password = document.getElementById("inputPasswordChanger")
    var newPassword = document.getElementById("inputNewPassword")

    if (password.value === newPassword.value) { // If passwords field are same

      await fetch(`${API_ADRESS}/resetPassword`, {
        method: "POST",
        headers:
        {
          "Authorization": token,
        },
        body: JSON.stringify({
          "UserPassword": password.value,
        })
      }).then(response => {

        if (response.status == 202) {

          alert("Successfully!")
          window.location.pathname = "/login"

        } else if (response.status == 401){

          alert("Something went wrong while authorization! Try to send recovery email again")
          window.location.pathname = "/reset"

        } else {

          alert("Something went wrong on a server! Try again later")
          window.location.reload()

        }

      })
    } else { // If passwords aren't same
      document.getElementById("newPasswordText").className = "newPassword incorrectChanger"
    }

  }

  return (
    <>
      <div id='rootChanger'>
        <div className="mainPageChanger">
          <h1>Write new password</h1>
          <div>
            <input type="password" className="inputPasswordChanger" id="inputPasswordChanger" />
            <i className="passwordChanger">Password</i>
          </div>
          <div>
            <input type="password" className="inputNewPassword" id="inputNewPassword" />
            <i className="newPassword" id="newPasswordText">Repeat Password</i>
          </div>
          <input type="button" className="ChangeChanger" value="Change" onClick={() => {resetPassword()}} />
        </div>
      </div>
    </>
  )
}
