/** @format */
"use client";

import { useState } from "react";
import { Nav } from "./ui/nav";

type Props = {};

import {
    ShoppingCart,
    LayoutDashboard,
    UsersRound,
    Settings,
    ChevronRight,
} from "lucide-react";
import { Button } from "./ui/button";

import { useWindowWidth } from "@react-hook/window-size";

export default function SideNavbar({}: Props) {
    const [isCollapsed, setIsCollapsed] = useState(false);

    const onlyWidth = useWindowWidth();
    const mobileWidth = onlyWidth < 768;

    function toggleSidebar() {
        setIsCollapsed(!isCollapsed);
    }

    return (
        <div className="relative w-[20rem] min-w-[200px] border-r px-3 pb-10 pt-10 ">
            <div className="mb-20 p-5">Logo</div>
            {/* {!mobileWidth && (
                <div className="absolute right-[-20px] top-7">
                    <Button
                        onClick={toggleSidebar}
                        variant="secondary"
                        className=" rounded-full p-2">
                        <ChevronRight />
                    </Button>
                </div>
            )} */}
            <Nav
                isCollapsed={mobileWidth ? true : isCollapsed}
                links={[
                    {
                        title: "Library",
                        href: "/library",
                        variant: "default",
                    },
                    {
                        title: "Parameters",
                        href: "/parameters",
                        variant: "ghost",
                    },
                    {
                        title: "CMS",
                        href: "/cms",
                        variant: "ghost",
                    },
                    {
                        title: "Settings",
                        href: "/settings",
                        variant: "ghost",
                    },
                ]}
            />
        </div>
    );
}
