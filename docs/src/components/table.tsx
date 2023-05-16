import React, {ReactNode} from "react";


type TableProps = {
    rows: [{
        category: ReactNode,
        rows: [{
            title: ReactNode
            details: ReactNode
        }]
    }]
}

export function Table({rows}: TableProps) {
    return (
        <table style={{width: '100%'}}>
            {
                rows.map(row => {
                    return (
                        <>
                            <tr style={{textAlign: 'center', width: '100%', fontWeight: 'bold'}}>
                                <td colSpan={2}> {row.category} </td>
                            </tr>
                            {row.rows.map(r => {
                                return (
                                    <tr>
                                        <td> {r.title} </td>
                                        <td> {r.details} </td>
                                    </tr>
                                )
                            })}
                        </>)
                })
            }
        </table>
    )
}