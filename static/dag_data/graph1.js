loadData(
    {
        name: 'graph1',
        nodes: [
            { id: 'table1', value: { label: '表1' } },
            { id: 'table2', value: { label: '表2' } },
            { id: 'table3', value: { label: '表3' } },
            { id: 'table4', value: { label: '表4' } },
            { id: 'table5', value: { label: '表5' } },
            { id: 'table6', value: { label: '表6' } },
            { id: 'table7', value: { label: '表7' } },
            { id: 'table8', value: { label: '表8' } },
            { id: 'table9', value: { label: '表9' } },
            { id: 'table10', value: { label: '表10' } },
            { id: 'table11', value: { label: '表11' } },
            { id: 'table12', value: { label: '表12' } },
            { id: 'table13', value: { label: '表13' } },
        ],
        links: [
            { u: 'table1', v: 'table2', value: { label: 'table1-table2' } },
            { u: 'table1', v: 'table3', value: { label: 'table1-table3' } },
            { u: 'table2', v: 'table4', value: { label: 'table2-table4' } },
            { u: 'table3', v: 'table4', value: { label: 'link4' } },
            { u: 'table4', v: 'table5', value: { label: 'link5' } },
            { u: 'table4', v: 'table6', value: { label: 'link5' } },
            { u: 'table4', v: 'table7', value: { label: 'link5' } },
            { u: 'table7', v: 'table8', value: { label: '测试8' } },
            { u: 'table7', v: 'table9', value: { label: '测试9' } },
            { u: 'table7', v: 'table10', value: { label: '测试10' } },
            { u: 'table10', v: 'table11', value: { label: '测试11' } },
            { u: 'table11', v: 'table12', value: { label: '测试12' } },
            { u: 'table12', v: 'table13', value: { label: '测试13' } },
        ]
    }
);