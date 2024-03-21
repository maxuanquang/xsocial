import { useContext, useEffect, useState } from "react";
import "./rightbar.css";
import { Link } from "react-router-dom";
import { AuthContext } from "../../context/AuthContext";
import { Add, Remove } from "@material-ui/icons";
import { makeRequest } from "../../axios";

export default function Rightbar({ user }) {
    const [userFollowers, setUserFollowers] = useState([]);
    const { user: currentUser, dispatch } = useContext(AuthContext);
    const [followed, setFollowed] = useState(false);

    useEffect(() => {
        const getFollowers = async () => {
            try {
                const response = await makeRequest.get(
                    `/friends/${user.user_id}/followers`
                );
                setFollowed(
                    response.data.followers_ids.includes(currentUser.user_id)
                );

                const followers = [];
                for (const id of response.data.followers_ids) {
                    try {
                        const userResponse = await makeRequest.get(
                            `/users/${id}`
                        );
                        followers.push(userResponse.data);
                    } catch (err) {
                        console.log(err);
                    }
                }
                setUserFollowers(followers);
            } catch (err) {}
        };
        getFollowers();
    }, [user]);

    const handleClick = async () => {
        try {
            if (followed) {
                await makeRequest.delete(`/friends/${user.user_id}`);
                dispatch({ type: "UNFOLLOW", payload: user.user_id });
            } else {
                await makeRequest.post(`/friends/${user.user_id}`);
                dispatch({ type: "FOLLOW", payload: user.user_id });
            }
            // setFollowed(!followed);
            window.location.reload(false);
        } catch (err) {
            console.log(err);
        }
    };

    const HomeRightbar = () => {
        return (
            <>
                <div className="birthdayContainer">
                    <img className="birthdayImg" src={"/gift.png"} alt="" />
                    <span className="birthdayText">
                        <b>Pola Foster</b> and <b>3 other friends</b> have a
                        birhday today.
                    </span>
                </div>
                <img className="rightbarAd" src={"/ad.png"} alt="" />
            </>
        );
    };

    const ProfileRightbar = () => {
        return (
            <>
                {user.user_name !== currentUser.user_name && (
                    <button
                        className="rightbarFollowButton"
                        onClick={handleClick}
                    >
                        {followed ? "Unfollow" : "Follow"}
                        {followed ? <Remove /> : <Add />}
                    </button>
                )}
                <h4 className="rightbarTitle">User followers</h4>
                <div className="rightbarFollowings">
                    {userFollowers.slice(0, 9).map((follower) => (
                        <Link
                            to={"/profile/" + follower.user_id}
                            style={{ textDecoration: "none" }}
                        >
                            <div className="rightbarFollowing">
                                <img
                                    src={
                                        follower.profile_picture
                                            ? follower.profile_picture
                                            : "/person/noAvatar.jpeg"
                                    }
                                    alt=""
                                    className="rightbarFollowingImg"
                                />
                                <span className="rightbarFollowingName">
                                    {follower.user_name}
                                </span>
                            </div>
                        </Link>
                    ))}
                </div>
            </>
        );
    };

    return (
        <div className="rightbar">
            <div className="rightbarWrapper">
                {user ? <ProfileRightbar /> : <HomeRightbar />}
            </div>
        </div>
    );
}
