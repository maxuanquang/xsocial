import Share from "../share/Share";
import Post from "../post/Post";
import "./feed.css";
import { useContext, useEffect, useState } from "react";
import { makeRequest } from "../../axios";
import { AuthContext } from "../../context/AuthContext";

export default function Feed(props) {
    const [posts, setPosts] = useState([]);
    const { user } = useContext(AuthContext);

    useEffect(() => {
        async function fetchPosts() {
            try {
                const response = props.user_id
                    ? await makeRequest.get(`/friends/${props.user_id}/posts`)
                    : await makeRequest.get(`/newsfeed`);

                const requestedPosts = [];
                for (const id of response.data.posts_ids) {
                    try {
                        const postResponse = await makeRequest.get(
                            `/posts/${id}`
                        );
                        requestedPosts.push(postResponse.data);
                    } catch (err) {
                        console.log(err);
                    }
                }
                setPosts(requestedPosts);
            } catch (err) {
                console.log(err);
            }
        }
        fetchPosts();
    }, [props.user_id]);

    return (
        <div className="feed">
            <div className="feedWrapper">
                {(!props.user_id || props.user_id === user.user_id) && (
                    <Share />
                )}
                {posts.map((p) => (
                    <Post key={p.post_id} post={p} />
                ))}
            </div>
        </div>
    );
}
