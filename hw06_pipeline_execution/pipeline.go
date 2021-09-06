package hw06

type (
	In  = <-chan interface{}
	Out = In
	Bi  = chan interface{}
)

type Stage func(in In) (out Out)

func ExecutePipeline(in In, done In, stages ...Stage) Out {
	if in == nil {
		ch := make(Bi)
		defer close(ch)
		return ch
	}
	for _, stage := range stages {
		ch := make(Bi)
		go func(ch Bi, in In) {
			defer close(ch)

			for {
				select {
				case <-done:
					return
				case val, ok := <-in:
					if !ok {
						return
					}
					ch <- val
				}
			}
		}(ch, in)
		in = stage(ch)
	}

	return in
}
