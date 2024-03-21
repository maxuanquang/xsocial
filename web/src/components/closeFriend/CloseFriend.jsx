import "./closeFriend.css"

export default function CloseFriend({user}) {
    return (
        <li className="sidebarFriend">
            <img className="sidebarFriendImg" src={user.profilePicture ? user.profilePicture : "/person/noAvatar.jpeg"} alt="" />
            <span className="sidebarFriendName">{user.username}</span>
        </li>
    )
}
