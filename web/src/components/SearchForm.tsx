import { Button, Card, Form, FormInstance, Input } from "antd";
import { useState } from "react";

import { useAPI } from "../lib/api";
import { YoutubeInfoJSON } from "../lib/types";

interface FormValues {
  target: string;
}

interface SearchFormProps {
  form: FormInstance<{ target: string }>;
  onMetadata: (metadata: YoutubeInfoJSON) => void;
}

export const SearchForm = (props: SearchFormProps) => {
  const [loading, setLoading] = useState(false);
  const api = useAPI();

  const onFinish = async (values: FormValues) => {
    setLoading(true);
    const metadata = await api.getMetadata(values.target);
    setLoading(false);
    props.onMetadata(metadata);
  };

  const clear = () => props.form.resetFields();

  return (
    <Card hoverable style={{ cursor: "auto" }}>
      <Form
        form={props.form}
        initialValues={{ target: "", audio_only: false }}
        onFinish={onFinish}
        style={{
          maxWidth: 800,
          width: "100%",
          marginLeft: "auto",
          marginRight: "auto",
        }}
        labelCol={{ span: 8 }}
        wrapperCol={{ span: 16 }}
      >
        <Form.Item
          label="Target"
          name="target"
          rules={[{ required: true }, { type: "url" }]}
        >
          <Input />
        </Form.Item>

        <Form.Item wrapperCol={{ offset: 8 }} style={{ marginBottom: 0 }}>
          <Button
            shape="round"
            type="primary"
            htmlType="submit"
            loading={loading}
          >
            Search
          </Button>

          <Button shape="round" onClick={clear} style={{ marginLeft: 8 }}>
            Clear
          </Button>
        </Form.Item>
      </Form>
    </Card>
  );
};
