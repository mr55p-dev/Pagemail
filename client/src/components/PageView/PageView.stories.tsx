import type { Meta, StoryObj } from "@storybook/react";

import { Page } from "./PageView.component";

const meta: Meta<typeof Page> = {
  component: Page,
  argTypes: {
    url: {
      control: "text",
    },
  },
};

export default meta;
type Story = StoryObj<typeof Page>;

export const Primary: Story = {
  render: ({ ...args }) => (
    <div style={{ maxWidth: "500px" }}>
      <Page {...args} />
    </div>
  ),
  args: {
    id: "123456",
    created: "2023-01-01T00:00:00",
	url: "https://www.example.com",
  },
};
