/*@flow*/
import React from 'react';
import FaCaretUp from 'react-icons/lib/fa/caret-up';


function Product(props: { id: number, title: string, vote: number, dispatch: Function }) {
    const { id, title, vote } = props;

    function handleVote() {
        props.dispatch({
            type: 'products/vote',
            id: id,
        });
    }

    return (
        <div>
            <div>
                {vote}
                <FaCaretUp onClick={handleVote} />
                {title}
            </div>
        </div>
    );

}



export default Product;
