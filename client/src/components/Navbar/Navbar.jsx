
// import React, { useEffect, useState } from "react";
// import { Link } from "react-router-dom";
// import { useSelector } from "react-redux";
// import { FaBars } from "react-icons/fa";
// import MobileMenu from "./MobileMenu";
// import NavItems from "./NavItems";
// import AvatarDropdown from "./AvatarDropDown";
// import Avatar from "../../assets/avatar/defaultAvatar.jpg"
// import { AiOutlineShoppingCart } from "react-icons/ai";

// const Navbar = () => {
//   const { totalItems } = useSelector((state) => state.cart);
//   const token = useSelector((state) => state.auth.token);
//   const [isScrolled, setIsScrolled] = useState(false);
//   const [isOpen, setIsOpen] = useState(false);
//   const { user } = useSelector((state) => state.profile);

//   const userProfileImage = user ? user.profile_picture : Avatar;

//   useEffect(() => {
//     const handleScroll = () => {
//       if (window.scrollY > 10) {
//         setIsScrolled(true);
//       } else {
//         setIsScrolled(false);
//       }
//     };

//     window.addEventListener("scroll", handleScroll);
//     return () => {
//       window.removeEventListener("scroll", handleScroll);
//     };
//   }, []);

//   const toggleMenu = () => {
//     setIsOpen(!isOpen);
//   };

//   return (
//     <nav
//       className={`fixed top-0 w-full transition-colors duration-300 z-50 ${isScrolled ? "bg-black" : "bg-transparent"
//         }`}
//     >
//       <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
//         <div className="grid grid-cols-2 md:grid-cols-3 h-16">
//           <div className="grid place-content-center text-2xl font-normal justify-start md:justify-center">
//             <Link
//               to="/"
//               className={`font-bold cursor-pointer ${isScrolled ? "text-white" : "text-black"
//                 }`}
//             >
//               Trivart
//             </Link>
//           </div>
//           <div className="flex items-center md:hidden justify-end">
//             <button
//               onClick={toggleMenu}
//               className="text-gray-500 hover:text-gray-700 "
//             >
//               <FaBars size={24} />
//             </button>
//           </div>
//           <div className="hidden md:flex lg:flex col-start-1 col-end-2 row-start-1 place-items-center">
//             <NavItems isScrolled={isScrolled} handleNavItemClick={toggleMenu} />
//           </div>
//           <div className="hidden md:flex items-center justify-center md:justify-end space-x-4">
//             {token ? (
//               <AvatarDropdown userProfileImage={userProfileImage} />
//             ) : (
//               <Link
//                 to="/signin"
//                 className={`${isScrolled ? "text-white" : "text-black"
//                   } hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium`}
//               >
//                 Sign In
//               </Link>
//             )}
//             <Link to="/cart" className="relative">
//               <AiOutlineShoppingCart
//                 className={`text-2xl ${isScrolled ? "text-white" : "text-black"
//                   }`}
//               />
//               {totalItems > 0 && (
//                 <span
//                   className={`absolute -bottom-2 -right-2 grid h-5 w-5 place-items-center overflow-hidden rounded-full text-center text-xs font-bold ${isScrolled
//                       ? "bg-red-600  text-white"
//                       : "bg-red-600 text-white"
//                     }`}
//                 >
//                   {totalItems}
//                 </span>
//               )}
//             </Link>
//           </div>
//         </div>
//         <MobileMenu
//           isOpen={isOpen}
//           toggleMenu={toggleMenu}
//           token={token}
//           userProfileImage={userProfileImage}
//           cartCount={totalItems}
//         />
//       </div>
//     </nav>
//   );
// };

// export default Navbar;

import React, { useEffect, useState } from "react";
import { Link, useLocation } from "react-router-dom";
import {useDispatch, useSelector} from "react-redux";
import { FaBars } from "react-icons/fa";
import MobileMenu from "./MobileMenu";
import NavItems from "./NavItems";
import AvatarDropdown from "./AvatarDropDown";
import Avatar from "../../assets/avatar/defaultAvatar.jpg";
import { AiOutlineShoppingCart } from "react-icons/ai";
import {
  Button,
  Drawer,
  DrawerBody,
  DrawerCloseButton,
  DrawerContent, DrawerFooter,
  DrawerHeader,
  DrawerOverlay, Input, useDisclosure
} from "@chakra-ui/react";
import { BsBagHeart } from "react-icons/bs";
import {removeFromFavorites} from "../../slices/favouriteSlice.jsx";


const Navbar = () => {
  const { totalItems } = useSelector((state) => state.cart);
  const token = useSelector((state) => state.auth.token);
  const { user } = useSelector((state) => state.profile);
  const [isModOpen, setModIsOpen] = useState(false);
  const [isVisible, setIsVisible] = useState(true);
  const [lastScrollY, setLastScrollY] = useState(0);
  const { isOpen, onOpen, onClose } = useDisclosure()
  const btnRef = React.useRef()

  const { favorites } = useSelector((state) => state.favourite);
  const dispatch=useDispatch();

  const userProfileImage = user ? user.profile_picture : Avatar;

  const location = useLocation();
  const isHomePage = location.pathname === "/";

  const handleRemoveFavorite = (productId) => {
    dispatch(removeFromFavorites(productId));
  };

  useEffect(() => {
    const handleScroll = () => {
      const currentScrollY = window.scrollY;

      if (currentScrollY > lastScrollY) {
        // User is scrolling down
        setIsVisible(false);
      } else {
        // User is scrolling up
        setIsVisible(true);
      }

      setLastScrollY(currentScrollY);
    };

    window.addEventListener("scroll", handleScroll);
    return () => {
      window.removeEventListener("scroll", handleScroll);
    };
  }, [lastScrollY]);

  const toggleMenu = () => {
    setModIsOpen(!isModOpen);
  };

  return (
    <>
      {/* Coupon Banner */}
      {isHomePage && (
        <div className="w-full bg-blue-500 text-white text-xs md:text-lg  text-center p-3 fixed top-0 z-30 font-Poppins">
          <p className={"font-Poppins"}>Use code <strong>SAVE10</strong> to get 10% off on your first purchase!</p>
        </div>
      )}

      {/* Navbar */}
      <nav
        className={`fixed w-full bg-black transition-transform duration-300 z-40 font-Poppins ${isVisible ? "translate-y-0" : "-translate-y-full"
          } ${isHomePage ? "mt-10 md:mt-[3.25rem]" : ""}`}  // Added margin to make space for the banner
      >
        <div className="max-w-7xl mx-auto px-4 sm:px-6 lg:px-8">
          <div className="grid grid-cols-2 md:grid-cols-3 h-16">
            <div className="grid place-content-center text-2xl font-normal justify-start md:justify-center">
              <Link to="/" className="font-bold cursor-pointer text-white">
                Trivart
              </Link>
            </div>
            <div className="flex gap-2 items-center md:hidden justify-end">
              <button ref={btnRef} onClick={onOpen} className={"pb-[4px]"}>
                <BsBagHeart color="white" size={24}/>
              </button>
              <button
                  onClick={toggleMenu}
                  className="text-gray-500 hover:text-gray-700"
              >
                <FaBars size={24}/>
              </button>

            </div>
            <div className="hidden md:flex lg:flex col-start-1 col-end-2 row-start-1 place-items-center">
            <NavItems handleNavItemClick={toggleMenu} />
            </div>
            <div className="hidden md:flex items-center justify-center md:justify-end space-x-4">
              {token ? (
                <AvatarDropdown userProfileImage={userProfileImage} />
              ) : (
                <Link
                  to="/signin"
                  className="text-white hover:bg-gray-700 hover:text-white px-3 py-2 rounded-md text-sm font-medium"
                >
                  Sign In
                </Link>
              )}
              <Link to="/cart" className="relative">
                <AiOutlineShoppingCart  className="text-2xl text-white" />
                {totalItems > 0 && (
                  <span className="absolute -bottom-2 -right-2 grid h-5 w-5 place-items-center overflow-hidden rounded-full text-center text-xs font-bold bg-red-600 text-white">
                    {totalItems}
                  </span>
                )}
              </Link>
              <button ref={btnRef}  onClick={onOpen}>
                <BsBagHeart color="white" size={25}/>
              </button>
              <Drawer
                  isOpen={isOpen}
                  placement='right'
                  onClose={onClose}
                  finalFocusRef={btnRef}
              >
                <DrawerOverlay />
                <DrawerContent>
                  <DrawerCloseButton />
                  <DrawerHeader>Favourite items</DrawerHeader>

                  <DrawerBody>

                    {favorites.length > 0 ? (
                        favorites.map((item) => (
                            <div
                                key={item.productId}
                                className="flex justify-between items-center mb-4"
                            >
                              <div>
                                <img
                                    src={item.image}
                                    alt={item.name}
                                    className="w-16 h-16 rounded-md"
                                />
                                <p>{item.name}</p>
                                <p>₹{item.price}</p>
                              </div>
                              <Button
                                  size="sm"
                                  colorScheme="red"
                                  onClick={() => handleRemoveFavorite(item.productId)}
                              >
                                Remove
                              </Button>
                            </div>
                        ))
                    ) : (
                        <p>No favorite items found</p>
                    )}

                  </DrawerBody>

                  <DrawerFooter>
                    <Button variant='outline' mr={3} onClick={onClose}>
                      Cancel
                    </Button>
                  </DrawerFooter>
                </DrawerContent>
              </Drawer>
            </div>
          </div>
        </div>
      </nav>

      <MobileMenu
        isOpen={isModOpen}
        toggleMenu={toggleMenu}
        token={token}
        userProfileImage={userProfileImage}
        cartCount={totalItems}
      />
    </>
  );
};

export default Navbar;
