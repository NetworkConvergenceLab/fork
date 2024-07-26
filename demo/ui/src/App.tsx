import '@mantine/core/styles.css';
import { Flex, MantineProvider, Image, Title, Button, NavLink } from '@mantine/core';
import testingScenarioImage from './assets/testing_scenario.png';
import operatorImage from './assets/operator_diagram.png';
import cpuLightImage from './assets/cpu_light_workload.png';
import cpuHeavyImage from './assets/cpu_heavy_workload.png';
import axios from 'axios';
import { useState } from 'react';

const baseURL = "http://192.168.2.125:4000"

function App() {
  const [singleClusterPlot, setSingleClusterPlot] = useState<string | null>(null)
  const [multiClusterPlot, setMultiClusterPlot] = useState<string | null>(null)

  const [singleClusterLoading, setSingleClusterLoading] = useState(false)
  const [multiClusterLoading, setMultiClusterLoading] = useState(false)

  const singleClusterClickHandler = async () => {
    setSingleClusterLoading(true)
    const res = await axios.post(`${baseURL}/attack/single`)
    if (res.status === 200) {
      setSingleClusterPlot(res.data['file'])
    }
    setSingleClusterLoading(false)
  }

  const multiClusterClickHandler = async () => {
    setMultiClusterLoading(true)
    const res = await axios.post(`${baseURL}/attack/multi`)
    if (res.status === 200) {
      setMultiClusterPlot(res.data['file'])
    }
    setMultiClusterLoading(false)
  }

  return (
    <MantineProvider>
      <Flex direction='column' gap='lg' py='md'>
        <Title order={2} ml='4rem'>FORK</Title>

        <Flex direction='row' w='100%'>
          <Flex direction='column' gap='xl' mr='6rem'>
            <Image radius='sm' h={250} w='auto' fit='contain' src={testingScenarioImage} />
            <Image radius='sm' h={400} w='auto' fit='contain' src={operatorImage} />

            <Image radius='sm' h={300} w='auto' fit='contain' src={cpuLightImage} />
            <Image radius='sm' h={300} w='auto' fit='contain' src={cpuHeavyImage} />
          </Flex>

          <Flex direction='column' mt='4rem'>
            <Flex direction='column' mb='3.5rem' gap='lg'>
              <Button
                loading={singleClusterLoading}
                onClick={singleClusterClickHandler}>
                Vegeta (Single cluster)
              </Button>

              {
                singleClusterPlot &&
                  <NavLink href={`${baseURL}/plot/${singleClusterPlot}`} target="_blank" rel="noopener" label='Show plot' />
              }
            </Flex>

            <Flex direction='column' gap='lg'>
              <Button loading={multiClusterLoading} onClick={multiClusterClickHandler}>
                Vegeta (Multi cluster)
              </Button>

              {
                multiClusterPlot &&
                  <NavLink href={`${baseURL}/plot/${multiClusterPlot}`} target="_blank" rel="noopener" label='Show plot' />
              }
            </Flex>
          </Flex>
        </Flex>
      </Flex>
    </MantineProvider>
  ) 
  
}

export default App
