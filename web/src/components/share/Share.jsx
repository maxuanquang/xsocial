import { useContext, useRef, useState } from "react";
import "./share.css";
import {
    PermMedia,
    Label,
    Room,
    EmojiEmotions,
    Cancel,
} from "@material-ui/icons";
import { AuthContext } from "../../context/AuthContext";
import { makeRequest } from "../../axios";

export default function Share() {
    const { user } = useContext(AuthContext);
    const content_text = useRef();
    const [file, setFile] = useState(null);
    const [isPosting, setIsPosting] = useState(false)

    const submitHandler = async (e) => {
        e.preventDefault();
        const newPost = {
            user_id: user.user_id,
            content_text: content_text.current.value,
            content_image_path: [],
        };

        if (file) {
            // get url
            const response = await makeRequest.get("/posts/url");
            const data = await response.data;
            const presignedURL = data.url;
            const filename = presignedURL.split("?")[0];

            // put to s3
            const res = await fetch(presignedURL, {
                method: "PUT",
                body: file,
                headers: {
                    "Content-Type": file.type,
                },
            });

            if (res.ok) {
                console.log("Image uploaded successfully");
                newPost.content_image_path.push(filename);
            } else {
                console.error("Failed to upload image:", response.statusText);
            }
        }

        // make request to api server
        try {
            await makeRequest.post("/posts", newPost);
        } catch (err) {
            console.log(err);
        }

        setIsPosting(true)
        setTimeout(() => window.location.reload(), 1000)
    };

    return (
        <div className="share">
            <div className="shareWrapper">
                <div className="shareTop">
                    <img
                        className="shareProfileImg"
                        src={
                            user.profile_picture
                                ? user.profile_picture
                                : "/person/noAvatar.jpeg"
                        }
                        alt=""
                    />
                    <input
                        className="shareInput"
                        placeholder={
                            "What's in your mind " + user.user_name + "?"
                        }
                        ref={content_text}
                    />
                </div>
                <hr className="shareHr" />
                {file && (
                    <div className="shareImgContainer">
                        <img
                            className="shareImg"
                            src={URL.createObjectURL(file)}
                            alt=""
                        />
                        <Cancel
                            className="shareImgCancel"
                            onClick={() => setFile(null)}
                        />
                    </div>
                )}
                <form className="shareBottom" onSubmit={submitHandler}>
                    <div className="shareOptions">
                        <label htmlFor="file" className="shareOption">
                            <PermMedia
                                className="shareOptionIcon"
                                htmlColor="tomato"
                            />
                            <span className="shareOptionText">
                                Photo or Video
                            </span>
                            <input
                                type="file"
                                id="file"
                                accept=".png,.jpeg,.jpg"
                                style={{ display: "none" }}
                                onChange={(e) => setFile(e.target.files[0])}
                            />
                        </label>
                        <div className="shareOption">
                            <Label
                                className="shareOptionIcon"
                                htmlColor="blue"
                            />
                            <span className="shareOptionText">Tag</span>
                        </div>
                        <div className="shareOption">
                            <Room
                                className="shareOptionIcon"
                                htmlColor="green"
                            />
                            <span className="shareOptionText">Location</span>
                        </div>
                        <div className="shareOption">
                            <EmojiEmotions
                                className="shareOptionIcon"
                                htmlColor="goldenrod"
                            />
                            <span className="shareOptionText">Feelings</span>
                        </div>
                    </div>
                    <button className="shareButton" type="submit">
                        {/* Share */}
                        {isPosting ? 'Posting...' : 'Share'}
                    </button>
                </form>
            </div>
        </div>
    );
}
