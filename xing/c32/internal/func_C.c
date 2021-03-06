/* Copyright (C) 2015-2020 김운하(UnHa Kim)  < unha.kim.ghts at gmail dot com >

이 파일은 GHTS의 일부입니다.

이 프로그램은 자유 소프트웨어입니다.
소프트웨어의 피양도자는 자유 소프트웨어 재단이 공표한 GNU LGPL 2.1판
규정에 따라 프로그램을 개작하거나 재배포할 수 있습니다.

이 프로그램은 유용하게 사용될 수 있으리라는 희망에서 배포되고 있지만,
특정한 목적에 적합하다거나, 이익을 안겨줄 수 있다는 묵시적인 보증을 포함한
어떠한 형태의 보증도 제공하지 않습니다.
보다 자세한 사항에 대해서는 GNU LGPL 2.1판을 참고하시기 바랍니다.
GNU LGPL 2.1판은 이 프로그램과 함께 제공됩니다.
만약, 이 문서가 누락되어 있다면 자유 소프트웨어 재단으로 문의하시기 바랍니다.
(자유 소프트웨어 재단 : Free Software Foundation, Inc.,
59 Temple Place - Suite 330, Boston, MA 02111-1307, USA)

Copyright (C) 2015-2020년 UnHa Kim (< unha.kim.ghts at gmail dot com >)

This file is part of GHTS.

GHTS is free software: you can redistribute it and/or modify
it under the terms of the GNU Lesser General Public License as published by
the Free Software Foundation, version 2.1 of the License.

GHTS is distributed in the hope that it will be useful,
but WITHOUT ANY WARRANTY; without even the implied warranty of
MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
GNU Lesser General Public License for more details.

You should have received a copy of the GNU Lesser General Public License
along with GHTS.  If not, see <http://www.gnu.org/licenses/>. */

#include <windows.h>
#include <stdio.h>

typedef void (__stdcall *ETK_GetAccountName)(const char* pszAcc, char* pszAccName, int nAccNameSize);
typedef void (__stdcall *ETK_GetAcctDetailName)(const char* pszAcc, char* pszAccName, int nAccNameSize);
typedef void (__stdcall *ETK_GetAcctNickName)(const char* pszAcc, char* pszAccName, int nAccNameSize);

HMODULE hModule = NULL;
ETK_GetAccountName etkGetAccountName = NULL;
ETK_GetAcctDetailName etkGetAcctDetailName = NULL;
ETK_GetAcctNickName etkGetAcctNickName = NULL;

void *XingApiDll() {
    if (hModule == NULL) {
        hModule = LoadLibrary(TEXT("XingAPI.dll"));
    }

    return (void *)hModule;
}

void GetAccountName(void *AccountNo, void *Buffer, int BufferSize) {
    memset(Buffer, 0, BufferSize);

    if (BufferSize < 41) {
        printf("C.GetAccountName() : BufferSize too small. %d\n", BufferSize);
        return;
    }

    if (etkGetAccountName == NULL) {
        FARPROC proc = GetProcAddress(hModule, "ETK_GetAccountName");
        if (proc == NULL) {
            printf("'ETK_GetAccountName' not found.\n");
            return;
        }

        etkGetAccountName = (ETK_GetAccountName)(proc);
    }

    etkGetAccountName((LPCTSTR)AccountNo, (LPSTR)Buffer, BufferSize);
}

void GetAcctDetailName(void *AccountNo, void *Buffer, int BufferSize) {
    memset(Buffer, 0, BufferSize);

    if (BufferSize < 41) {
        printf("C.GetAcctDetailName() : BufferSize too small. %d\n", BufferSize);
        return;
    }

    if (etkGetAcctDetailName == NULL) {
        FARPROC proc = GetProcAddress(hModule, "ETK_GetAcctDetailName");
        if (proc == NULL) {
            printf("'ETK_GetAcctDetailName' not found.\n");
            return;
        }

        etkGetAcctDetailName = (ETK_GetAcctDetailName)(proc);
    }

    etkGetAcctDetailName((LPCTSTR)AccountNo, (LPSTR)Buffer, BufferSize);
}

void GetAcctNickName(void *AccountNo, void *Buffer, int BufferSize) {
    memset(Buffer, 0, BufferSize);

    if (BufferSize < 41) {
        printf("C.GetAcctNickName() : BufferSize too small. %d\n", BufferSize);
        return;
    }

    if (etkGetAcctNickName == NULL) {
        FARPROC proc = GetProcAddress(hModule, "ETK_GetAcctNickname");
        if (proc == NULL) {
            printf("'ETK_GetAcctNickName' not found.\n");
            return;
        }

        etkGetAcctNickName = (ETK_GetAcctNickName)(proc);
    }

    etkGetAcctNickName((LPCTSTR)AccountNo, (LPSTR)Buffer, BufferSize);
}