import "./post.css";
import { MoreVert } from "@material-ui/icons";
import { useState, useEffect, useContext } from "react";
import axios from "axios";
import { Link } from "react-router-dom";
import { format } from "timeago.js";
import { AuthContext } from "../../context/AuthContext";
import { makeRequest } from "../../axios";

export default function Post({ post }) {
    const PF = process.env.REACT_APP_PUBLIC_FOLDER;

    // const [like, setLike] = useState(post.likes.length);
    // const [isLiked, setIsLiked] = useState(false);
    // const { user: currentUser } = useContext(AuthContext);

    // const likeHandler = () => {
    //     try {
    //         axios.post(`/posts/${post.post_id}/likes`);
    //     } catch (err) {}
    //     setLike(isLiked ? like - 1 : like + 1);
    //     setIsLiked(!isLiked);
    // };

    const [postUser, setPostUser] = useState({});
    useEffect(() => {
        async function fetchUser() {
            const response = await makeRequest.get(`/users/${post.user_id}`);
            setPostUser(response.data);
        }
        fetchUser();
    }, [post.user_id]); // Only re-run the effect if post.userId changes

    return (
        <div className="post">
            <div className="postWrapper">
                <div className="postTop">
                    <div className="postTopLeft">
                        <Link to={`/profile/${postUser.user_id}`}>
                            <img
                                className="postProfileImg"
                                src={
                                    postUser.profile_picture ? postUser.profile_picture : "/person/noAvatar.jpeg"
                                }
                                alt=""
                            />
                        </Link>
                        <span className="postUsername">
                            {postUser.user_name}
                        </span>
                        <span className="postDate">
                            {format(post.created_at)}
                            {" " + post.created_at}
                        </span>
                    </div>
                    <div className="postTopRight">
                        <MoreVert />
                    </div>
                </div>
                <div className="postCenter">
                    <span className="postText">{post?.content_text}</span>
                    {post.content_image_path && (
                        <img
                            className="postImg"
                            src={post.content_image_path[0]}
                            alt=""
                        />
                    )}
                </div>
                <div className="postBottom">
                    <div className="postBottomLeft">
                        <img
                            className="likeIcon"
                            src={"/like.png"}
                            // onClick={likeHandler}
                            alt=""
                        />
                        <img
                            className="likeIcon"
                            src={"/heart.png"}
                            // onClick={likeHandler}
                            alt=""
                        />
                        <span className="postLikeCounter">
                            {post.users_liked ? post.users_liked.length : 0}{" "}
                            people like it
                        </span>
                    </div>
                    <div className="postBottomRight">
                        <span className="postCommentText">
                            {post.comments ? post.comments.length : 0} comments
                        </span>
                    </div>
                </div>
            </div>
        </div>
    );
}
