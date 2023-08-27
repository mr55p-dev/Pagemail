import type { Meta, StoryObj } from "@storybook/react";

import { Page } from "./PageView.component";
import { Grid } from "@mui/joy";

const pagePreviewData = JSON.parse(
  `{
    "title": "Using Service Workers - Web APIs | MDN",
    "description": "This article provides information on getting started with service workers, including basic architecture, registering a service worker, the installation and activation process for a new service worker, updating your service worker, cache control and custom responses, all in the context of a simple app with offline functionality.",
    "url": "https://developer.mozilla.org/en-US/docs/Web/API/Service_Worker_API/Using_Service_Workers"
}`
);

const meta: Meta<typeof Page> = {
  component: Page,
  parameters: {
    mockData: [
      {
		url: `http://127.0.0.1:8090/api/preview?target=https%3A%2F%2Fwww.example.com`,
        method: "GET",
        status: 200,
        response: pagePreviewData,
      },
    ],
  },
};

export default meta;
type Story = StoryObj<typeof Page>;

export const Primary: Story = {
  render: ({ ...args }) => (
    <Grid container>
      <Page {...args} />
    </Grid>
  ),
  args: {
    id: "123456",
    created: "2023-01-01T00:00:00",
    url: "https://www.example.com",
  },
};
