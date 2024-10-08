import React, { useCallback, useState } from "react";
import { RiDeleteBin6Line } from "react-icons/ri";
import PropTypes from "prop-types";
import { apiConnector } from "../../services/apiConnector.jsx";
import { Image, useToast } from "@chakra-ui/react";
import { Link } from "react-router-dom";
import { debounce } from "../../utils/debounce.js";


export const CartItem = ({ item, refreshCart }) => {
  const [quantity, setQuantity] = useState(item.quantity);
  const toast = useToast();

  /*

  const handleQuantityChange = async (change) => {
    // const newQuantity = Math.max(1, quantity + amount);
    setQuantity((prevQuantity) => {
      const newQuantity = prevQuantity + change;
      if (newQuantity < 5 && newQuantity > 0) {
        return newQuantity;
      }
      return prevQuantity;
    });

    try {
      const response = await apiConnector(
        "POST",
        "update-item",
        { cart_item_id: item.id, quantity },
        null,
        null,
        true
      );
      console.log("Quantity change", response);
      toast({
        title: response?.data?.msg || "item updated",
        status: "success",
        duration: 2500,
        isClosable: true,
      });
      refreshCart();
    } catch (e) {
      toast({
        title: e?.data?.msg || "error updating",
        status: "error",
        duration: 2500,
        isClosable: true,
      });
    }
    // Here you might want to update the quantity in the backend or global state
  };





  */

    // };
    const debouncedUpdateQuantity = useCallback(
      debounce(async (cartItemId, newQuantity) => {
        try {
          const response = await apiConnector(
            "POST",
            "update-item",
            { cart_item_id: cartItemId, quantity: newQuantity },
            null,
            null,
            true
          );
          console.log("Quantity change", response);
          toast({
            title: response?.data?.msg || "item updated",
            status: "success",
            duration: 2500,
            isClosable: true,
          });
          refreshCart();
        } catch (e) {
          toast({
            title: e?.data?.msg || "error updating",
            status: "error",
            duration: 2500,
            isClosable: true,
          });
        }
      }, 500), // 500ms debounce delay
      []
    );
  
    const handleQuantityChange = (change) => {
      setQuantity((prevQuantity) => {
        const newQuantity = prevQuantity + change;
        if (newQuantity < 999 && newQuantity > 0) {
          debouncedUpdateQuantity(item.id, newQuantity);
          return newQuantity;
        }
  
        if (newQuantity > 5) {
          toast({
            title: "Quantity Limit Reached",
            description: "You cannot add more than 5 of this item to your cart.",
            status: "warning",
            duration: 2500,
            isClosable: true,
          });
        }
        return prevQuantity;
      });
    };
  
    const handleItemDelete = async (id) => {
      // Handle item deletion here, possibly with a dispatch or API call
      try {
        const response = await apiConnector(
          "DELETE",
          `/remove-item/${id}`,
          null,
          null,
          null,
          true
        );
        toast({
          title: response?.data?.msg || "item removed",
          status: "success",
          duration: 2500,
          isClosable: true,
        });
        console.log("itemdelete", response);
        refreshCart();
      } catch (e) {
        toast({
          title: e?.data?.msg || "error removing",
          status: "error",
          duration: 2500,
          isClosable: true,
        });
      }
    };
  
    return (

        <div className="w-full flex flex-col md:flex-row md:justify-between md:items-center p-5 border-b-2">
            <div className={"flex-col md:flex-row flex  gap-4"}>
                <div>
                    <Image src={item.product_image} className={"w-[30%] md:w-[10rem] rounded-sm"}/>
                </div>
                <div className="flex md:items-center">
                    <div className="">
                        <h3 className="text-lg md:text-xl font-semibold">{item.product_name}</h3>
                        <p className="text-grey-600 line-through">₹{item.price / item.quantity}</p>
                        <p className="text-gray-600">₹{item.after_discount_total_price / item.quantity} <span className="text-green-600">({Math.floor(((item.price-item.after_discount_total_price)/item.price)*100)}% off)</span></p> 
                        <Link to={`/product/${item.product_id}`}>
                        <button className="text-blue-500 hover:underline mt-2 md:mt-0">More Details</button>
                        </Link>
                    </div>
                </div>
            </div>
            <div className="flex justify-between items-center w-full md:w-auto mt-3">
                <div className="flex items-center md:mr-8">
                    <button
                        onClick={() => handleQuantityChange(-1)}
                        className="px-2 py-1 bg-gray-200 rounded"
                    >
                        -
                    </button>
                    <span className="mx-4">{quantity}</span>
                    <button
                        onClick={() => handleQuantityChange(1)}
                        className="px-2 py-1 bg-gray-200 rounded"
                    >
                        +
                    </button>
                </div>
                <p className="text-lg font-semibold md:mb-0 mr-0 md:mr-8">₹{item.after_discount_total_price}</p>
                <button
                    onClick={() =>handleItemDelete(item.id)}
                    className="text-red-500 text-xl"
                >
                    <RiDeleteBin6Line />
                </button>
            </div>
            </div>

            )}
  {/* // Debounced API call for updating the quantity */}
  
 

CartItem.propTypes = {
  item: PropTypes.shape({
    id: PropTypes.string.isRequired,
    product_name: PropTypes.string.isRequired,
    price: PropTypes.number.isRequired,
    quantity: PropTypes.number.isRequired,
    after_discount_total_price: PropTypes.number.isRequired,
    product_id: PropTypes.string.isRequired,
    product_image: PropTypes.string.isRequired,
  }).isRequired,
};

export default CartItem;
