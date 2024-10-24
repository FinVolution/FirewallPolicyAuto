import { useMenusDispatch } from "@/context/menu";
import { addInternet, getFirewalls, getInternetList } from "@/server";
import { App, Button, Col, Flex, Form, Input, Modal, Radio, Row, Select, Space, Table, theme, Typography } from "antd";
import { ColumnsType } from "antd/lib/table";
import _ from "lodash";
import qs from "qs";
import React, { Fragment, useEffect, useState } from "react";
import { useSearchParams } from 'react-router-dom'


const { Paragraph } = Typography


const Policy = () => {
  const { message } = App.useApp()
  const [searchParams, setSearchParams] = useSearchParams()
  const [data, setData] = useState<any[]>()
  const [search, setSearch] = useState<any>(() => {
    const obj: any = {}
    searchParams.forEach((value, key) => {
      obj[key] = value
    })
    if (Object.keys(obj).length) {
      return obj
    }
    return {}
  })
  const [page, setPage] = useState<number>(1)
  const [pageSize, setPageSize] = useState<number>(15)
  const [total, setTotal] = useState<number>(0)
  const [lodaing, setLoading] = useState<boolean>(false)
  const [form] = Form.useForm()
  const { token } = theme.useToken()
  const dispatch = useMenusDispatch()
  const [firewalls, setFirewalls] = useState<any[]>([])
  const [virtualZones, setVirtualZones] = useState<any[]>([])
  const [formVirtualZones, setFormVirtualZones] = useState<any[]>([])
  const [open, setOpen] = useState(false)
  const [flag, setFlag] = useState<boolean>(false)
  const [confirmLoading, setConfirmLoading] = useState<boolean>(false)
  useEffect(() => {
    dispatch({
      type: 'changed',
      selectedKeys: 'policy',
      items: [
        {
          title: '网络策略'
        },
        {
          title: '策略查询'
        }
      ]
    })
  }, [])

  useEffect(() => {
    getFirewalls().then(res => {
      const address = searchParams.get('address')
      if (res.data.code === 0) {
        if (address) {
          const data = res.data.data.find((v: any) => v.address === address)
          setVirtualZones(data.virtualZone)
        }
        setFirewalls(res.data.data)
      } else {
        setFirewalls([])
      }
    })
  }, [])


  const tagstyle: React.CSSProperties = {
    background: token.colorFillQuaternary,
    color: token.colorText,
    borderColor: token.colorBorder,
    borderRadius: token.borderRadiusOuter,
    border: `1px solid ${token.colorBorder}`,
    paddingLeft: token.paddingXS,
    paddingRight: token.paddingXS,
    marginRight: token.marginXS,
    whiteSpace: 'nowrap'
  }

  useEffect(() => {
    form.setFieldsValue(search)
  }, [form])


  useEffect(() => {
    const params = {
      page,
      pageSize,
      ...search,
    }
    setLoading(true)
    getInternetList(params).then(res => {
      if (res.data.code === 0) {
        setData(res.data.data.policyList)
        setTotal(res.data.data.total)
      } else {
        setData([])
      }
      setLoading(false)
    })

  }, [search, page, pageSize, flag])

  const columns: ColumnsType<any> = [
    {
      title: 'ID',
      dataIndex: 'id',
    },
    {
      title: '名称',
      dataIndex: 'name',
    },
    {
      title: '动作',
      dataIndex: 'action',
      render: (v) => {
        return v === 1 ? <span style={{ color: 'red' }} >拒绝</span> : '允许'
      }
    },
    {
      title: '状态',
      dataIndex: 'enable',
      render: (v) => {
        return v ? "启用" : <span style={{ color: 'red' }} >禁用</span>
      }
    },
    {

      title: '源地址',
      dataIndex: 'srcAddress',
      width: 400,
      render: (val, obj, index) => {
        const row = 1
        return <Paragraph ellipsis={{ rows: row, expandable: 'collapsible' }} style={{ margin: 0, width: '100%', }} >
          {
            val?.map((v: string, i: number) => {
              return <span key={'net' + i} style={{ ...tagstyle, }} >{v}&nbsp;</span>
            })
          }
        </Paragraph>
      }
    },
    {
      title: '目的地址',
      width: 400,
      dataIndex: 'destAddress',
      render: (val, obj, index) => {
        const row = 1
        return <Paragraph ellipsis={{ rows: row, expandable: 'collapsible' }} style={{ margin: 0, width: '100%', }} >
          {
            val?.map((v: string, i: number) => {
              return <span key={'net' + i} style={{ ...tagstyle, }} >{v}&nbsp;</span>
            })
          }
        </Paragraph>
      }
    },
    {
      title: '源安全域',
      dataIndex: 'srcZone',
      width: 100,
      render: (val) => {
        if (!val) return '-'
        return val.map((v: string, i: number) => {
          return <div key={v}>{v}</div>
        })
      }
    },
    {
      title: '目的安全域',
      dataIndex: 'destZone',
      width: 100,
      render: (val, obj, index) => {
        const row = 1
        return <Paragraph ellipsis={{ rows: row, expandable: 'collapsible' }} style={{ margin: 0 }}>
          {
            val.map((v: string, i: number) => {
              return <span key={i} >{v}{(i === val.length - 1) ? null : '，'}</span>
            })
          }
        </Paragraph>
      }
    },
    {
      title: '协议/端口',
      dataIndex: 'servicePort',
      width: 200,
      render: (val, obj, index) => {
        const row = 1
        return <Paragraph ellipsis={{ rows: row, expandable: 'collapsible' }} style={{ margin: 0, width: '100%', }} >
          {
            val?.map((v: string, i: number) => {
              return <span key={'net' + i} style={{ ...tagstyle, }} >{v}&nbsp;</span>
            })
          }
        </Paragraph>
      }
    },
    {
      title: '防火墙名称',
      dataIndex: 'firewallName',
    },
  ]

  const finish = (e: any) => {
    setPage(1)
    setSearch(e)
    const str = qs.stringify(e)
    const obj: any = qs.parse(str)
    setSearchParams(obj)
  }

  const style: React.CSSProperties = {
    margin: 0
  }
  const submit = (values: any) => {
    setConfirmLoading(true)
    addInternet(values).then(res => {
      if (res.data.code === 0) {
        message.success('操作成功')
        setOpen(false)
        setFlag(!flag)
      } else {
        message.error(res.data.msg)
      }
    }).finally(() => {
      setConfirmLoading(false)
    })
  }
  return <Fragment>
    <Flex vertical gap={8} style={{ transform: 'translateY(-24px)', }} >
      <div style={{ position: 'sticky', top: 0, zIndex: 2, background: 'rgba(255, 255, 255, 0.9)', paddingTop: '24px', paddingBottom: '12px' }} >
        <Form onFinish={finish} form={form} >
          <Row gutter={[24, 8]} style={{ width: '100%' }} >
            <Col span={6}>
              <Form.Item label="防火墙地址" name="address" style={style}>
                <Select
                  options={firewalls}
                  fieldNames={{ label: 'name', value: 'address' }}
                  allowClear
                  optionRender={(val) => {
                    return <span>{val.value}-{val.label}</span>
                  }}
                  labelRender={(val) => {
                    return val.value
                  }}
                  onChange={(e, option) => {
                    if (option) {
                      setVirtualZones(option.virtualZone)
                    } else {
                      setVirtualZones([])
                    }
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item label="虚拟墙" name="virtualZone" style={style}>
                <Select
                  options={virtualZones}
                  fieldNames={{ label: 'name', value: 'code' }}
                  optionRender={(val) => {
                    return <span>{val.value}-{val.label}</span>
                  }}
                  labelRender={(val) => {
                    return val.value
                  }}
                />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item label="源地址" name="srcAddr" style={style}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item label="目的地址" name="dstAddr" style={style}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item label="源区域" name="srcZone" style={style}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item label="目的区域" name="destZone" style={style}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={6}>
              <Form.Item label="端口" name="port" style={style}>
                <Input />
              </Form.Item>
            </Col>
            <Col span={6}  >
              <Form.Item wrapperCol={{ span: 24 }} style={{ ...style, textAlign: 'right' }}  >
                <Space size="large" >
                  <Button onClick={() => setOpen(true)} >新建策略</Button>
                  <Button htmlType="reset">重置</Button>
                  <Button type="primary" htmlType="submit">搜索</Button>
                </Space>
              </Form.Item>
            </Col>
          </Row>
        </Form>
      </div>

      <Table
        bordered
        columns={columns}
        dataSource={data}
        loading={lodaing}
        size="small"
        rowKey="id"
        scroll={{
          x: 'max-content'
        }}
        pagination={{
          pageSize,
          current: page,
          total,
          pageSizeOptions: ['10', '15', '20', '30', '40', '50'],
          showQuickJumper: true,
          showTotal: (total) => `共${total}条`,
          onChange: (page, pageSize) => {
            setPage(page)
            setPageSize(pageSize)
          }
        }}
      />
    </Flex>
    <Modal title="新建策略" open={open} onCancel={() => setOpen(false)} destroyOnClose footer={null} >
      <Form preserve={false} onFinish={submit} >
        <Form.Item label="标题" name="title" rules={[{ required: true }]} >
          <Input />
        </Form.Item>
        <Form.Item label="源区域" name="srcZone" >
          <Input />
        </Form.Item>
        <Form.Item label="目的区域" name="destZone" >
          <Input />
        </Form.Item>
        <Form.Item label="源地址" name="srcAddr" rules={[{ required: true }]} >
          <Select
            mode="tags"
            notFoundContent={null}
          ></Select>
        </Form.Item>
        <Form.Item label="目的地址" name="destAddr" rules={[{ required: true }]} >
          <Select
            mode="tags"
            notFoundContent={null}
          ></Select>
        </Form.Item>
        <Form.Item label="端口" name="service" rules={[{ required: true }]} >
          <Select
            mode="tags"
            notFoundContent={null}
          ></Select>
        </Form.Item>
        <Form.Item label="访问" name="action" rules={[{ required: true }]} >
          <Radio.Group>
            <Radio value={1} >拒绝</Radio>
            <Radio value={2} >允许</Radio>
          </Radio.Group>
        </Form.Item>
        <Form.Item label="防火墙地址" name="firewallAddress" rules={[{ required: true }]} >
          <Select
            options={firewalls}
            fieldNames={{ label: 'name', value: 'address' }}
            allowClear
            optionRender={(val) => {
              return <span>{val.value}-{val.label}</span>
            }}
            labelRender={(val) => {
              return val.value
            }}
            onChange={(e, option) => {
              if (option) {
                setFormVirtualZones(option.virtualZone)
              } else {
                setFormVirtualZones([])
              }
            }}
          />
        </Form.Item>
        <Form.Item label="虚拟墙" name="virtualZone"  >
          <Select
            options={formVirtualZones}
            fieldNames={{ label: 'name', value: 'code' }}
            optionRender={(val) => {
              return <span>{val.value}-{val.label}</span>
            }}
            labelRender={(val) => {
              return val.value
            }}
          />
        </Form.Item>
        <Form.Item>
          <Space style={{ float: 'right' }} >
            <Button onClick={() => setOpen(false)} loading={confirmLoading}>取消</Button>
            <Button type="primary" htmlType="submit" loading={confirmLoading} >保存</Button>
          </Space>
        </Form.Item>
      </Form>
    </Modal>
  </Fragment>
}

export default Policy

