'use server';

import axios from "axios";
import { createSearchParamsFromClient } from "next/dist/server/request/search-params";
import { NextRequest, NextResponse } from "next/server";
import { cookies } from 'next/headers'


export async function POST(request: NextRequest) {
  const formData = await request.json();
  const sessionCookie = request.cookies.get('brick-session-id')
  var response;
  try {
    const resp = await axios.post(
        'http://localhost:8080/logout',
        {
            headers: {
                Cookie: `${sessionCookie?.name}=${sessionCookie?.value}`
            },
            withCredentials: true
        }
    );
    response = NextResponse.json({
      status: "Ok"
    });
  } catch (err) {
    if (axios.isAxiosError(err)) {
      response = NextResponse.json(err.response?.data || {status: 'Error'});
    } else {
        response = NextResponse.json({status: 'Error'});
    }
  }
  if (sessionCookie?.name){
    console.log('cookie found', sessionCookie.name, sessionCookie.value)
    response.cookies.delete(sessionCookie?.name)
    }
  return response
}