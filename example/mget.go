package example

type rediscli struct {
}

func (c *rediscli) MGET(keys []string) ([][]byte, error) {
	return nil, nil
}

func MGET(wrappers ...ProtoWrapper) error {

	cli := &rediscli{}

	keys := make([]string, 0, len(wrappers))
	for _, v := range wrappers {
		keys = append(keys, v.GetKey())
	}

	result, err := cli.MGET(keys)
	if err != nil {
		return err
	}

	for i, v := range result {
		err = wrappers[i].Unmarshal(v)
		if err != nil {
			return err
		}
	}

	return nil
}
