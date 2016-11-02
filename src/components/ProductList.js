/*@flow*/
import React from 'react';
import Product from './Product';

function ProductList(props: { data: Array < { ID: number, title: string, vote: number } > , dispatch: Function }) {
    return (
        <div >
        {
          props.data.map(product => <Product
            key={product.ID}
            id={product.ID}
            title={product.title}
            vote={product.vote}
            dispatch={props.dispatch}
          />)
        }
    </div>
    );
}


export default ProductList;
