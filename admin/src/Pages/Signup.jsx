import React,{useState} from 'react'
import { ArrowRight } from 'lucide-react'
import {Link, useNavigate} from "react-router-dom";
import Cover1 from "../assets/Cover/cover-1.jpg"
import {useDispatch, useSelector} from "react-redux";
import {signupUser} from "../redux/authSlice.js";
import {
    useToast
} from '@chakra-ui/react';



export function Signup() {
    const [Username, setUsername] = useState('');
    const [Password, setPassword] = useState('');
    const toast=useToast({position: 'top',})
    const dispatch = useDispatch();
    const navigate = useNavigate();
    const signupStatus=useSelector((state)=>state.auth.status)

    const handleSignup =async () => {
      try{
          const response=await dispatch(signupUser({ Username, Password }))
          if (response.payload.error) {
              toast({
                  title: 'Signup failed.',
                  description: response?.payload?.error || 'An error occurred.',
                  status: 'error',
                  duration: 3000,
                  isClosable: true,
              });

              }else{

              toast({
                  title: 'Signup successful.',
                  description: 'Redirecting to login...',
                  status: 'success',
                  duration: 3000,
                  isClosable: true,
              });
              navigate("/login")

          }

      }catch (err){
          toast({
              title: 'Signup failed.',
              description: err.message || 'An error occurred.',
              status: 'error',
              duration: 3000,
              isClosable: true,
          });
          console.log(err);
      }


      }


    return (
        <section>
            <div className="flex flex-row w-screen h-screen items-center justify-center pr-5 ">
                <div
                    className="flex w-[50%] flex-col items-center justify-around px-4 py-10 sm:px-6 sm:py-16 lg:px-16 lg:py-24">
                    <div className="w-full ">
                        <h2 className="text-3xl font-bold leading-tight text-black sm:text-4xl">Sign up</h2>
                        <p className="mt-2 text-base text-gray-600">
                            Already have an account?{' '}
                            <Link to={"/login"}>
                            <span

                                className="font-medium text-black transition-all duration-200 hover:underline"
                            >
                                Sign In
                            </span>
                            </Link>
                        </p>
                        <form action="#" method="POST" className="mt-8 w-[80%]">
                            <div className="space-y-5">
                                <div>


                                </div>
                                <div>
                                    <label htmlFor="email" className="text-base font-medium text-gray-900">
                                        {' '}
                                        Email address{' '}
                                    </label>
                                    <div className="mt-2">
                                        <input
                                            className="flex h-10 w-full rounded-md border border-gray-300 bg-transparent px-3 py-2 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400 focus:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-50"
                                            type="email"
                                            placeholder="Email"
                                            id="email"
                                            value={Username}
                                            onChange={(e) => setUsername(e.target.value)}
                                        ></input>
                                    </div>
                                </div>
                                <div>
                                    <div className="flex items-center justify-between">
                                        <label htmlFor="password" className="text-base font-medium text-gray-900">
                                            {' '}
                                            Password{' '}
                                        </label>
                                    </div>
                                    <div className="mt-2">
                                        <input
                                            className="flex h-10 w-full rounded-md border border-gray-300 bg-transparent px-3 py-2 text-sm placeholder:text-gray-400 focus:outline-none focus:ring-1 focus:ring-gray-400 focus:ring-offset-1 disabled:cursor-not-allowed disabled:opacity-50"
                                            type="password"
                                            placeholder="Password"
                                            id="password"
                                            value={Password}
                                            onChange={(e) => setPassword(e.target.value)}
                                        ></input>
                                    </div>
                                </div>
                                <div>
                                    <button
                                        type="button"
                                        className="inline-flex w-full items-center justify-center rounded-md bg-black px-3.5 py-2.5 font-semibold leading-7 text-white hover:bg-black/80"
                                        onClick={handleSignup}
                                        disabled={signupStatus === 'loading'}

                                    >
                                        Create Account <ArrowRight className="ml-2" size={16}/>
                                    </button>
                                </div>
                            </div>
                        </form>

                    </div>
                </div>
                <div className="h-[80%] w-[50%] relative flex items-center bg-red-300 justify-center rounded-md  ">
                    <img
                        className="h-full w-full rounded-md object-cover"
                        src={Cover1}
                        alt="Cover Image for Signup Page"
                    />
                    <div className="absolute inset-0 bg-black opacity-65 transition-opacity duration-300 rounded-md "></div>
                    <div className="absolute inset-0 flex items-center justify-center rounded-md ">
                        <span className="text-white text-3xl font-[700] font-Jost">ForeverKnots Admin Page</span>
                    </div>
                </div>

            </div>
        </section>
    );
}
