/**
 * eslint-disable @next/next/no-img-element
 *
 * @format
 */

/** @format */
"use client";

import { Button } from "@/components/ui/button";

import React from "react";
import PageTitle from "@/components/PageTitle";

type Props = {};

export default function UsersPage({}: Props) {
    return (
        <div>
            <PageTitle title="Library" className="pb-10 text-5xl font-normal" />
            <div className="flex flex-row gap-5  w-full">
                <Button variant="outline" className="w-40 py-12 ">
                    Upload Video
                </Button>
                <Button variant="outline" className="w-40 py-12">
                    Create Folder
                </Button>
            </div>
        </div>
    );
}
