'use server';

import axios from "axios";
import { NextRequest, NextResponse } from "next/server";

export async function POST(
    request: NextRequest,
) {
    const formData = await request.json();
    try {
        await axios.post('http://localhost:8080/admin', formData)
        return NextResponse.json({
            status: "Ok"
        })
    } catch (err) {
        if (axios.isAxiosError(err)) {
            return NextResponse.json(err.response?.data || {status: 'Error'})
        } else {
            return NextResponse.json({status: 'Error'});
        }
    }
}