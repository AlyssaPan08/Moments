import React, { useState, useEffect } from "react";
import PropTypes from "prop-types";
import { message } from "antd";
import { DeleteOutlined } from "@ant-design/icons";
import axios from "axios";
import Gallery from "react-grid-gallery";
import { BASE_URL, TOKEN_KEY } from "../constants";

const captionStyle = {
    backgroundColor: "rgba(0, 0, 0, 0.6)",
    maxHeight: "240px",
    overflow: "hidden",
    position: "absolute",
    bottom: "0",
    width: "100%",
    color: "white",
    padding: "2px",
    fontSize: "90%"
};

const wrapperStyle = {
    display: "block",
    minHeight: "1px",
    width: "100%",
    border: "1px solid #ddd",
    overflow: "auto"
};

function PhotoGallery(props) {
    const [images, setImages] = useState(props.images);
    const [curImgIdx, setCurImgIdx] = useState(0);

    const imageArr = images.map((image) => {
        return {
            ...image,
            customOverlay: (
                <div style={captionStyle}>
                    <div>{`${image.user}: ${image.caption}`}</div>
                </div>
            )
        };
    });

    const onDeleteImage = () => {
        //step 1: get the image id to be deleted
        //step 2: filter the image from the image arr
        //step 3: send delete request to the server
        //step 4: analyze the response from the server
        //  case 1: success -> update the state: image
        //  case 2: fail -> inform user
        if (window.confirm(`Are you sure you want to delete this image?`)) {
            const curImg = images[curImgIdx];
            const newImageArr = images.filter((img, index) => index !== curImgIdx);
            console.log("delete image ", newImageArr);
            const opt = {
                method: "DELETE",
                url: `${BASE_URL}/post/${curImg.postId}`,
                headers: {
                    Authorization: `Bearer ${localStorage.getItem(TOKEN_KEY)}`
                }
            };

            axios(opt)
                .then((res) => {
                    console.log("delete result -> ", res);
                    // case1: success
                    if (res.status === 200) {
                        // step1: set state
                        setImages(newImageArr);
                    }
                })
                .catch((err) => {
                    // case2: fail
                    message.error("Fetch posts failed!");
                    console.log("fetch posts failed: ", err.message);
                });
        }
    };

    const onCurrentImageChange = (index) => {
        console.log("curIdx ", index);
        setCurImgIdx(index);
    };

    useEffect(() => {
        setImages(props.images);
    }, [props.images]);

    return (
        <div style={wrapperStyle}>
            <Gallery
                images={imageArr}
                enableImageSelection={false}
                backdropClosesModal={true}
                currentImageWillChange={onCurrentImageChange}
                customControls={[
                    <button
                        style={{ marginTop: "10px", marginLeft: "5px" }}
                        key="deleteImage"
                        type="primary"
                        icon={<DeleteOutlined />}
                        size="small"
                        onClick={onDeleteImage}
                    >
                        Delete Image
                    </button>
                ]}
            />
        </div>
    );
}

//校验prop type: images should be array of ...
PhotoGallery.propTypes = {
    images: PropTypes.arrayOf(
        PropTypes.shape({
            user: PropTypes.string.isRequired,
            caption: PropTypes.string.isRequired,
            src: PropTypes.string.isRequired,
            thumbnail: PropTypes.string.isRequired,
            thumbnailWidth: PropTypes.number.isRequired,
            thumbnailHeight: PropTypes.number.isRequired
        })
    ).isRequired
};

export default PhotoGallery;