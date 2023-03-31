package com.transport.irc_messenger

import android.os.Bundle
import androidx.fragment.app.Fragment
import android.view.LayoutInflater
import android.view.View
import android.view.ViewGroup
import androidx.navigation.fragment.findNavController
import com.transport.irc_messenger.databinding.FragmentFirstBinding

/**
 * A simple [Fragment] subclass as the default destination in the navigation.
 */
class FirstFragment : Fragment() {

    private var _binding: FragmentFirstBinding? = null

    // This property is only valid between onCreateView and
    // onDestroyView.
    private val binding get() = _binding!!

    override fun onCreateView(
        inflater: LayoutInflater, container: ViewGroup?,
        savedInstanceState: Bundle?
    ): View? {

        _binding = FragmentFirstBinding.inflate(inflater, container, false)
        return binding.root

    }

    override fun onViewCreated(view: View, savedInstanceState: Bundle?) {
        super.onViewCreated(view, savedInstanceState)

        val bruh = irc_transport.Irc_transport.newIrcTransport("bruh_test_kotel")
        val msg = irc_transport.Message()
        msg.setUserId("aaa")
        msg.setText("BRUH BRUH BRUH")
        msg.setMessageId(42)
        msg.setParentId(42)
        msg.setTimestamp(102093)
        binding.buttonFirst.setOnClickListener {
            bruh.sendMessages("#test_kotel_channel", msg)
        }
    }

    override fun onDestroyView() {
        super.onDestroyView()
        _binding = null
    }
}