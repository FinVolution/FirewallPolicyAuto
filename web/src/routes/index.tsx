import {
  createBrowserRouter,
  redirect,
} from "react-router-dom";
import LayoutPage from "../layout";
import Policy from "@/pages/Policy";

const router = createBrowserRouter([
  {
    path: "/",
    element: <LayoutPage />,
    loader: ({request}) => {
      const url = new URL(request.url);
      if(url.pathname === '/'){
        return redirect("/policy")
      }else{
        return ''
      }
    },
    children: [
      {
        path: "policy",
        element: <Policy />
      },
    ]
  }
]);

export default router